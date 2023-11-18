package ohio_ring_buffer2

import (
	"fmt"
	"io"
)

const RingBufferDefaultSize = 4096

type RingBuffer struct {
	dataLength              int
	readLength, writeLength int64
	rIdx, wIdx              int
	data                    []byte
}

func MakeRingBuffer(n int) *RingBuffer {
	if n == 0 {
		n = RingBufferDefaultSize
	}

	return &RingBuffer{
		data: make([]byte, n),
	}
}

// func (rb *RingBuffer) String() string {

// }

func (rb *RingBuffer) Reset() {
	rb.dataLength = 0
	rb.readLength = 0
	rb.writeLength = 0
	rb.rIdx = 0
	rb.wIdx = 0

	for i := range rb.data {
		rb.data[i] = 0
	}
}

func (rb *RingBuffer) PeekWriteable() (rs RingSlice) {
	if rb.Len() == len(rb.data) {
		return
	}

	rs.start = rb.writeLength

	if rb.wIdx < rb.rIdx {
		rs.data[0] = rb.data[rb.wIdx:rb.rIdx]
	} else {
		rs.data[1] = rb.data[:rb.rIdx]
		rs.data[0] = rb.data[rb.wIdx:]
	}

	wCap := rs.Len()

	if wCap > len(rb.data) {
		panic(
			fmt.Sprintf(
				"wcap was %d but buffer len was %d",
				wCap,
				len(rb.data),
			),
		)
	}

	return
}

func (rb *RingBuffer) PeekReadable() (rs RingSlice) {
	if rb.Len() == 0 {
		return
	}

	rs.start = rb.readLength

	if rb.rIdx < rb.wIdx {
		rs.data[0] = rb.data[rb.rIdx:rb.wIdx]
	} else {
		rs.data[1] = rb.data[:rb.wIdx]
		rs.data[0] = rb.data[rb.rIdx:]
	}

	rCap := rs.Len()

	if rCap > rb.Len() {
		panic(
			fmt.Sprintf(
				"rcap was %d but buffer len was %d and len was %d and n was %d and r was %d and w was %d",
				rCap,
				len(rb.data),
				rb.Len(),
				rb.dataLength,
				rb.rIdx,
				rb.wIdx,
			),
		)
	}

	return
}

func (rb *RingBuffer) Cap() int {
	return len(rb.data)
}

func (rb *RingBuffer) Write(p []byte) (n int, err error) {
	if rb.Len() == len(rb.data) {
		err = io.EOF
		return
	}

	rs := rb.PeekWriteable()

	var n1 int

	n1 = copy(rs.data[0], p)
	rb.wIdx += n1
	n += n1
	rb.dataLength += n1
	rb.writeLength += int64(n1)

	if rb.Len() == len(rb.data) {
		err = io.EOF
		return
	}

	if n == len(p) {
		return
	}

	n1 = copy(rs.data[1], p[n:])
	n += n1
	rb.dataLength += n1
	rb.writeLength += int64(n1)

	if n1 > 0 {
		rb.wIdx = n1
	}

	if rb.Len() == len(rb.data) {
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

	n1 = copy(p, rs.data[0])
	rb.rIdx += n1
	rb.dataLength -= n1
	rb.readLength += int64(n1)
	n += n1

	if rb.Len() == 0 {
		err = io.EOF
		return
	}

	if n == len(p) {
		return
	}

	n1 = copy(p[n:], rs.data[1])
	n += n1
	rb.dataLength -= n1
	rb.readLength += int64(n1)

	if n1 > 0 {
		rb.rIdx = n1
	}

	if rb.Len() == 0 {
		err = io.EOF
		return
	}

	return
}

func (rb *RingBuffer) FillWith(r io.Reader) (n int, err error) {
	if rb.Len() == len(rb.data) {
		err = io.EOF
		return
	}

	rs := rb.PeekWriteable()

	var n1 int

	n1, err = r.Read(rs.First())
	rb.wIdx += n1
	n += n1
	rb.dataLength += n1
	rb.writeLength += int64(n1)

	if err != nil {
		return
	}

	if rb.Len() == len(rb.data) {
		err = io.EOF
		return
	}

	n1, err = r.Read(rs.Second())
	n += n1
	rb.dataLength += n1
	rb.writeLength += int64(n1)

	if n1 > 0 {
		rb.wIdx = n1
	}

	if rb.Len() == len(rb.data) {
		err = io.EOF
		return
	}

	return
}

func (rb *RingBuffer) advanceRead(n int) {
	rb.rIdx += n
	rb.readLength += int64(n)

	if rb.rIdx > len(rb.data) {
		rb.rIdx -= len(rb.data)
	}

	rb.dataLength -= n
}

func (rb *RingBuffer) FindNext(ff FindFunc) (offset, length int, partial bool) {
	rs := rb.PeekReadable()
	return rs.findAnywhere(ff)
}

func (rb *RingBuffer) FindAnywhere(ff FindFunc) (offset, length int, partial bool) {
	return rb.PeekReadable().findAnywhere(ff)
}

func (rb *RingBuffer) FindFromStart(ff FindFunc) (length int, partial bool) {
	return rb.PeekReadable().findFromStart(ff)
}

func (rb *RingBuffer) FindFromStartAndAdvance(ff FindFunc) (length int, partial bool) {
	length, partial = rb.PeekReadable().findFromStart(ff)

	if !partial {
		rb.advanceRead(length)
	}

	return
}

func (rb *RingBuffer) Len() int {
	if rb.dataLength > len(rb.data) {
		panic("length is greater than buffer")
	}

	return rb.dataLength
}
