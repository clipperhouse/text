package segment

import (
	"errors"
	"fmt"
)

// SegmentFunc is the primitive, stateless unit for segmentation, and is
// where your implementation logic will live. It takes a byte slice and returns
// the boundaries of the first segment (token) in that slice.
//
// start is the index of the first byte of the first segment.
//
// end is the index of the first byte *following* the first segment.
// It's intended to be used for slicing, e.g. data[start:end].
//
// To indicate that there are no valid segments, return end = 0. (For example,
// a whitespace splitter operating on data that is only whitespace.)
//
// SegmentFunc is similar in spirit to bufio.SplitFunc, but offers more
// granular position information.
type SegmentFunc func(data []byte, atEOF bool) (start int, end int, err error)

// Segmenter is an iterator for byte arrays. See the New() and Next() funcs.
type Segmenter struct {
	f    SegmentFunc
	data []byte

	start, end int
	previous   *stack

	err error
}

var _ Forward = &Segmenter{}
var _ Bidirectional = &Segmenter{}

// New creates a new segmenter given a SegmentFunc. To use the new segmenter,
// call SetText() and then iterate using Next()
func New(f SegmentFunc) *Segmenter {
	return &Segmenter{
		f:        f,
		previous: &stack{},
	}
}

// SetText sets the text for the segmenter to operate on, and resets
// all state (i.e. the current position)
func (sc *Segmenter) SetText(data []byte) {
	sc.data = data

	sc.start = 0
	sc.end = 0
	sc.previous.clear()
	sc.err = nil
}

// Next advances the Segmenter to the next segment. It returns false when there
// are no remaining segments, or an error occurred.
//	text := []byte("This is an example.")
//
//	segment := whitespace.NewSegmenter(text)
//	for segment.Next() {
//		fmt.Printf("%q\n", segment.Bytes())
//	}
//	if err := segment.Err(); err != nil {
//		log.Fatal(err)
//	}
func (seg *Segmenter) Next() bool {
	if seg.end == len(seg.data) {
		return false
	}

	start, end, err := seg.f(seg.data[seg.end:], true)
	seg.err = err

	if seg.err != nil {
		return false
	}

	if start > end {
		seg.err = fmt.Errorf("the start of the next segment (%d) is greater than the end (%d); this is likely a bug in the SegmentFunc",
			start, end)
		return false
	}

	if end == 0 || start == end { // i.e. no segment
		return false
	}

	if end > len(seg.data) {
		seg.err = fmt.Errorf("the end of the next segment (%d) exceeds the length of the text (%d); this is likely a bug in the SegmentFunc",
			end, len(seg.data))
		return false
	}

	if seg.end != 0 {
		seg.previous.push(seg.start, seg.end)
	}

	seg.start = seg.end + start
	seg.end = seg.end + end

	return true
}

// Previous moves the segmenter to the previous segment. It returns false when there
// are no remaining previous segments, or an error occurred.
func (seg *Segmenter) Previous() bool {
	if start, end, exists := seg.previous.pop(); exists {
		seg.start = start
		seg.end = end
		return true
	}

	return false
}

// Start is the index of the first byte of the current segment
func (seg *Segmenter) Start() int {
	return seg.start
}

// End is the index of the first byte following the end of the current segment.
// The bytes of the current segment are data[segment.Start():segment.End()]
func (seg *Segmenter) End() int {
	return seg.end
}

// Err indicates an error occured when calling Next() or Previous(). Next and
// Previous will return false when an error occurs.
func (seg *Segmenter) Err() error {
	return seg.err
}

// Bytes returns the current segment
func (seg *Segmenter) Bytes() []byte {
	return seg.data[seg.start:seg.end]
}

var ErrIncompleteRune = errors.New("incomplete rune")
var ErrIncompleteToken = errors.New("incomplete token")

// ToSplitFunc maps a SegmentFunc to a bufio.SplitFunc, as a convenience
func ToSplitFunc(f SegmentFunc) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(data) == 0 {
			return 0, nil, nil
		}

		start, end, err := f(data, atEOF)

		if err != nil && !atEOF {
			if errors.Is(err, ErrIncompleteRune) {
				// Rune extends past current data, request more
				return 0, nil, nil
			}

			if errors.Is(err, ErrIncompleteToken) {
				// Token extends past current data, request more
				return 0, nil, nil
			}
		}

		return end, data[start:end], err
	}
}
