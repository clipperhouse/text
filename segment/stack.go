package segment

type span struct {
	start, end int
}

type stack struct {
	spans []span
}

func (st *stack) len() int {
	return len(st.spans)
}

func (st *stack) push(start, end int) {
	item := span{start, end}
	st.spans = append(st.spans, item)
}

func (st *stack) pop() (start, send int, exists bool) {
	if len(st.spans) == 0 {
		return 0, 0, false
	}

	last := len(st.spans) - 1
	span := st.spans[last]
	st.spans = st.spans[:last]

	return span.start, span.end, true
}

func (st *stack) clear() {
	st.spans = st.spans[:0]
}
