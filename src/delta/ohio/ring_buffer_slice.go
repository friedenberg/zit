package ohio

import (
	"fmt"
)

type RingSlice [2][]byte

func (rs RingSlice) String() string {
	return fmt.Sprintf("first: %q, second: %q", rs.First(), rs.Second())
}

func (rs RingSlice) First() []byte {
	return rs[0]
}

func (rs RingSlice) Second() []byte {
	return rs[1]
}

func (rs RingSlice) IsEmpty() bool {
	return rs.Len() == 0
}

func (rs RingSlice) Len() int {
	return len(rs.First()) + len(rs.Second())
}

func (rs RingSlice) Find(ff FindFunc) (offset int, eof bool) {
	if rs.Len() == 0 {
		return
	}

	return ff(rs)
}

// rs is the data to search. Negative offset means not found. 0 or positive
// offset means found at that index. Partial means the sequence is not complete.
type FindFunc func(rs RingSlice) (offset int, partial bool)

func FindBoundary(m []byte) FindFunc {
	return func(rs RingSlice) (offset int, partial bool) {
		offset = -1

		if len(m) == 0 {
			return
		}

		i := 0
		j := 0
		lastWasMatch := false

		for _, v := range rs.First() {
			if m[i] != v {
				lastWasMatch = false
				i = 0
			} else {
				lastWasMatch = true
				i++

				if i == len(m) {
					break
				}
			}

			j++
		}

		if i < len(m) {
			for _, v := range rs.Second() {
				if m[i] != v {
					lastWasMatch = false
					i = 0
				} else {
					lastWasMatch = true
					i++

					if i == len(m) {
						break
					}
				}

				j++
			}
		}

		switch {
		case i == len(m) && !lastWasMatch:
			panic("last was not match but match was read completely")

		case i > len(m) && lastWasMatch:
			panic("last was match but i is greater than length of m")
			// log.Debug().Printf("no boundary???")
			// offset = j

		case i == len(m) && lastWasMatch:
			if j == 0 {
				j += 1
			}

			offset = j - i

		case i < len(m)-1 && lastWasMatch:
			offset = j - i
			partial = true

		default:
		}

		return
	}
}

// example sequence:
// [text "quoted text []"]
// [text "quoted text []"]
// ^     ^            ^
