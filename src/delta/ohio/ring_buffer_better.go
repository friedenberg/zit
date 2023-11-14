package ohio

type MatchType int

const (
	MatchTypeNone     = MatchType(iota)
	MatchTypePartial  = MatchType(iota)
	MatchTypeComplete = MatchType(iota)
)

type MatchFunc func(data RingSlice, atEOF bool) (int, MatchType)

// func (rb *RingBuffer) PeekMatchAdvanceButSuperGOod(
// 	boundary []byte,
// 	atEOF bool,
// ) (n int, advance bool) {
// 	if rb.Len() < len(boundary) {
// 		advance = true
// 	}

// 	mf := func(rs RingSlice, eof bool) (n int, mt MatchType) {
// 		r := rb.r

// 		for _, v := range rs.First() {
// 			if n == len(boundary) {
// 				return MatchTypeComplete
// 			}

// 			if boundary[n] != v {
// 				return MatchTypeNone
// 			}

// 			n++
// 			r++
// 		}

// 		for _, v := range rs.Second() {
// 			if n == len(boundary) {
// 				return MatchTypeComplete
// 			}

// 			if boundary[n] != v {
// 				return MatchTypeNone
// 			}

// 			n++
// 			r++
// 		}

// 		if n == len(boundary) {
// 			return MatchTypeComplete
// 		} else {
// 			return MatchTypePartial
// 		}
// 	}

// 	n = rb.peekMatchAdvanceButBetter(mf, advance, atEOF)

// 	return
// }

// // TODO: modify `m` to be a function instead of a literal slice of bytes
// func (rb *RingBuffer) peekMatchAdvanceButBetter(
// 	mf MatchFunc,
// 	shouldAdvance,
// 	atEOF bool,
// ) (n int) {
// 	rs := rb.PeekReadable()

// 	n, mt := mf(rs, atEOF)

// 	switch mt {
// 	case MatchTypeNone:
// 		return

// 	case MatchTypePartial:
// 	case MatchTypeComplete:
// 	default:
// 		panic(fmt.Sprintf("unsupported match type: %#v", mt))
// 	}

// 	return
// }
