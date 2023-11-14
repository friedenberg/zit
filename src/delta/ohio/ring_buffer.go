package ohio

import (
	"fmt"
	"io"
)

const ringBufferDefaultSize = 4096

type RingBuffer struct {
	n, r, w int
	buffer  []byte
}

func MakeRingBuffer(n int) *RingBuffer {
	if n == 0 {
		n = ringBufferDefaultSize
	}

	return &RingBuffer{
		buffer: make([]byte, n),
	}
}

// func (rb *RingBuffer) String() string {

// }

func (rb *RingBuffer) Reset() {
	rb.n = 0
	rb.r = 0
	rb.w = 0

	for i := range rb.buffer {
		rb.buffer[i] = 0
	}
}

func (rb *RingBuffer) PeekWriteable() (rs RingSlice) {
	if rb.Len() == len(rb.buffer) {
		return
	}

	if rb.w < rb.r {
		rs[0] = rb.buffer[rb.w:rb.r]
	} else {
		rs[1] = rb.buffer[:rb.r]
		rs[0] = rb.buffer[rb.w:]
	}

	wCap := rs.Len()

	if wCap > len(rb.buffer) {
		panic(
			fmt.Sprintf(
				"wcap was %d but buffer len was %d",
				wCap,
				len(rb.buffer),
			),
		)
	}

	return
}

func (rb *RingBuffer) PeekReadable() (rs RingSlice) {
	if rb.Len() == 0 {
		return
	}

	if rb.r < rb.w {
		rs[0] = rb.buffer[rb.r:rb.w]
	} else {
		rs[1] = rb.buffer[:rb.w]
		rs[0] = rb.buffer[rb.r:]
	}

	rCap := rs.Len()

	if rCap > rb.Len() {
		panic(
			fmt.Sprintf(
				"rcap was %d but buffer len was %d and len was %d and n was %d and r was %d and w was %d",
				rCap,
				len(rb.buffer),
				rb.Len(),
				rb.n,
				rb.r,
				rb.w,
			),
		)
	}

	return
}

func (rb *RingBuffer) Cap() int {
	return len(rb.buffer)
}

func (rb *RingBuffer) Write(p []byte) (n int, err error) {
	if rb.Len() == len(rb.buffer) {
		err = io.EOF
		return
	}

	rs := rb.PeekWriteable()

	var n1 int

	n1 = copy(rs[0], p)
	rb.w += n1
	n += n1
	rb.n += n1

	if rb.Len() == len(rb.buffer) {
		err = io.EOF
		return
	}

	if n == len(p) {
		return
	}

	n1 = copy(rs[1], p[n:])
	n += n1
	rb.n += n1

	if n1 > 0 {
		rb.w = n1
	}

	if rb.Len() == len(rb.buffer) {
		err = io.EOF
		return
	}

	return
}

func (rb *RingBuffer) Read(p []byte) (n int, err error) {
	if rb.Len() == 0 {
		err = io.EOF
		return
	}

	rs := rb.PeekReadable()

	var n1 int

	n1 = copy(p, rs[0])
	rb.r += n1
	rb.n -= n1
	n += n1

	if rb.Len() == 0 {
		err = io.EOF
		return
	}

	if n == len(p) {
		return
	}

	n1 = copy(p[n:], rs[1])
	n += n1
	rb.n -= n1

	if n1 > 0 {
		rb.r = n1
	}

	if rb.Len() == 0 {
		err = io.EOF
		return
	}

	return
}

func (rb *RingBuffer) FillWith(r io.Reader) (n int, err error) {
	if rb.Len() == len(rb.buffer) {
		err = io.EOF
		return
	}

	rs := rb.PeekWriteable()

	var n1 int

	n1, err = r.Read(rs[0])
	rb.w += n1
	n += n1
	rb.n += n1

	if err != nil {
		return
	}

	if rb.Len() == len(rb.buffer) {
		err = io.EOF
		return
	}

	n1, err = r.Read(rs[1])
	n += n1
	rb.n += n1

	if n1 > 0 {
		rb.w = n1
	}

	if rb.Len() == len(rb.buffer) {
		err = io.EOF
		return
	}

	return
}

func (rb *RingBuffer) advance(n int) {
	rb.r += n

	if rb.r > len(rb.buffer) {
		rb.r -= len(rb.buffer)
	}

	rb.n -= n
}

func (rb *RingBuffer) FindFromStartAndAdvance(m []byte) (length int, partial bool) {
	length, partial = rb.PeekReadable().FindFromStart(FindBoundary(m))

	if !partial {
		rb.advance(length)
	}

	return
}

func (rb *RingBuffer) Len() int {
	if rb.n > len(rb.buffer) {
		panic("length is greater than buffer")
	}

	return rb.n
}
