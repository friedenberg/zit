package ohio_ring_buffer

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

func (rs RingSlice) FindFromStart(ff FindFunc) (length int, partial bool) {
	var offset int

	offset, length, partial = ff(rs)

	if offset > 0 {
		length = 0
		partial = false
	}

	return
}

func (rs RingSlice) FindAnywhere(ff FindFunc) (offset, length int, partial bool) {
	if rs.Len() == 0 {
		return
	}

	return ff(rs)
}

// rs is the data to search. Negative offset means not found. 0 or positive
// offset means found at that index. Partial means the sequence is not complete.
type FindFunc func(rs RingSlice) (offset, length int, partial bool)

func FindBoundary(boundary []byte) FindFunc {
	return func(rs RingSlice) (offset, length int, partial bool) {
		offset = -1

		if len(boundary) == 0 {
			return
		}

		j := 0
		lastWasMatch := false

		for _, v := range rs.First() {
			if boundary[length] != v {
				lastWasMatch = false
				length = 0
			} else {
				lastWasMatch = true
				length++

				if length == len(boundary) {
					break
				}
			}

			j++
		}

		if length < len(boundary) {
			for _, v := range rs.Second() {
				if boundary[length] != v {
					lastWasMatch = false
					length = 0
				} else {
					lastWasMatch = true
					length++

					if length == len(boundary) {
						break
					}
				}

				j++
			}
		}

		switch {
		case length == len(boundary) && !lastWasMatch:
			panic("last was not match but match was read completely")

		case length > len(boundary) && lastWasMatch:
			panic("last was match but i is greater than length of m")

		case length == len(boundary) && lastWasMatch:
			if j == 0 {
				j += 1
			}

			offset = j - length

		case length < len(boundary) && lastWasMatch:
			offset = j - length
			partial = true

		default:
		}

		return
	}
}
