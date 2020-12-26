package segment_test

import (
	"fmt"
	"testing"

	"golang.org/x/text/segment"
	"golang.org/x/text/segment/whitespace"
)

// See words/segmenter_test.go for a real test

// TODO: test Segmenter logic independent of SegmentFunc implementation

func TestNextPrevious(t *testing.T) {
	segment := segment.New(whitespace.SegmentFunc)
	segment.SetText([]byte("hi   how are you!!  \nand more\r"))

	for segment.Next() {
		fmt.Printf("%q\n", segment.Bytes())
	}

	if err := segment.Err(); err != nil {
		t.Error(err)
	}

	segment.SetText([]byte("Let's try previous"))
	if segment.Next() {
		fmt.Printf("%q\n", segment.Bytes())
	}
	if segment.Next() {
		fmt.Printf("%q\n", segment.Bytes())
	}
	if segment.Next() {
		fmt.Printf("%q\n", segment.Bytes())
	}
	if segment.Previous() {
		fmt.Printf("%q\n", segment.Bytes())
	}
	if segment.Previous() {
		fmt.Printf("%q\n", segment.Bytes())
	}
	if segment.Previous() {
		fmt.Printf("%q\n", segment.Bytes())
	}
	if segment.Previous() {
		fmt.Printf("%q\n", segment.Bytes())
	}
	if segment.Next() {
		fmt.Printf("%q\n", segment.Bytes())
	}
}
