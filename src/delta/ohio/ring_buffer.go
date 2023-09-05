package ohio

import (
	"fmt"
	"io"
	"math"
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

func (rb *RingBuffer) PeekWriteable() (first, second []byte) {
	if rb.Len() == len(rb.buffer) {
		return
	}

	if rb.w < rb.r {
		first = rb.buffer[rb.w:rb.r]
	} else {
		second = rb.buffer[:rb.r]
		first = rb.buffer[rb.w:]
	}

	wCap := len(first) + len(second)

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

func (rb *RingBuffer) PeekReadable() (first, second []byte) {
	if rb.Len() == 0 {
		return
	}

	if rb.r < rb.w {
		first = rb.buffer[rb.r:rb.w]
	} else {
		second = rb.buffer[:rb.w]
		first = rb.buffer[rb.r:]
	}

	rCap := len(first) + len(second)

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
	}

	first, second := rb.PeekWriteable()

	var n1 int

	n1 = copy(first, p)
	rb.w += n1
	n += n1
	rb.n += n1

	if rb.Len() == len(rb.buffer) {
		err = io.EOF
	}

	if n == len(p) {
		return
	}

	n1 = copy(second, p[n:])
	n += n1
	rb.n += n1

	if n1 > 0 {
		rb.w = n1
	}

	if rb.Len() == len(rb.buffer) {
		err = io.EOF
	}

	return
}

func (rb *RingBuffer) Read(p []byte) (n int, err error) {
	if rb.Len() == 0 {
		err = io.EOF
	}

	first, second := rb.PeekReadable()

	var n1 int

	n1 = copy(p, first)
	rb.r += n1
	rb.n -= n1
	n += n1

	if rb.Len() == 0 {
		err = io.EOF
	}

	if n == len(p) {
		return
	}

	n1 = copy(p[n:], second)
	n += n1
	rb.n -= n1

	if n1 > 0 {
		rb.r = n1
	}

	if rb.Len() == 0 {
		err = io.EOF
	}

	return
}

func (rb *RingBuffer) ReadFromSmall(r io.Reader) (n int, err error) {
	var n1 int64
	n1, err = rb.ReadFrom(r)

	if n1 > math.MaxInt {
		err = ErrReadFromSmallOverflow
		return
	}

	n = int(n1)

	return
}

func (rb *RingBuffer) ReadFrom(r io.Reader) (n int64, err error) {
	if rb.Len() == len(rb.buffer) {
		err = io.EOF
		return
	}

	first, second := rb.PeekWriteable()

	var n1 int

	n1, err = r.Read(first)
	rb.w += n1
	n += int64(n1)
	rb.n += n1
	if err != nil {
		return
	}

	if rb.Len() == len(rb.buffer) {
		err = io.EOF
		return
	}

	n1, err = r.Read(second)
	n += int64(n1)
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

func (rb *RingBuffer) Find(m []byte) (offset int, eof bool) {
	offset = -1

	if len(m) == 0 {
		return
	}

	if rb.Len() == 0 {
		return
	}

	first, second := rb.PeekReadable()
	i := 0
	j := 0
	lastWasMatch := false

	for _, v := range first {
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
		for _, v := range second {
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

	case i == len(m) && lastWasMatch:
		offset = j - i

	case i < len(m)-1 && lastWasMatch:
		offset = j - i
		eof = true

	default:
	}

	return
}

func (rb *RingBuffer) PeekMatchAdvance(m []byte) (n int, eof bool) {
	advance := true

	if rb.Len() < len(m) {
		advance = false
		eof = true
	}

	n = rb.peekMatchAdvance(m, advance)
	return
}

func (rb *RingBuffer) PeekMatch(m []byte) (n int) {
	return rb.peekMatchAdvance(m, false)
}

func (rb *RingBuffer) peekMatchAdvance(m []byte, advance bool) (n int) {
	r := rb.r

	first, second := rb.PeekReadable()

	for _, v := range first {
		if n == len(m) {
			break
		}

		if m[n] != v {
			return
		}

		n++
		r++
	}

	// if len(second) > 0 {
	// 	r = 0
	// }

	for _, v := range second {
		if n == len(m) {
			break
		}

		if m[n] != v {
			return
		}

		n++
		r++
	}

	if advance && n == len(m) {
		rb.r += n

		if rb.r > len(rb.buffer) {
			rb.r -= len(rb.buffer)
		}

		rb.n -= n
	}

	return
}

func (rb *RingBuffer) Len() int {
	if rb.n > len(rb.buffer) {
		panic("length is greater than buffer")
	}

	return rb.n
}
