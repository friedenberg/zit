package ohio_ring_buffer2

import (
	"fmt"
)

type RingSlice struct {
	start int64
	data  [2][]byte
}

func (rs RingSlice) String() string {
	return fmt.Sprintf("first: %q, second: %q", rs.First(), rs.Second())
}

func (rs RingSlice) Start() int64 {
	return rs.start
}

func (rs RingSlice) First() []byte {
	return rs.data[0]
}

func (rs RingSlice) Second() []byte {
	return rs.data[1]
}

func (rs RingSlice) IsEmpty() bool {
	return rs.Len() == 0
}

func (rs RingSlice) Len() int {
	return len(rs.First()) + len(rs.Second())
}

func (rs RingSlice) findFromStart(ff FindFunc) (length int, partial bool) {
	var offset int

	offset, length, partial = ff(rs)

	if offset > 0 {
		length = 0
		partial = false
	}

	return
}

func (rs RingSlice) findAnywhere(ff FindFunc) (offset, length int, partial bool) {
	if rs.Len() == 0 {
		return
	}

	return ff(rs)
}
