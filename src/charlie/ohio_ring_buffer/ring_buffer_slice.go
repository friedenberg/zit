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