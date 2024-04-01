package catgut

import (
	"fmt"
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/collections"
)

const RingBufferDefaultSize = 4096

func MakeRingBuffer(r io.Reader, n int) *RingBuffer {
	if n == 0 {
		n = RingBufferDefaultSize
	}

	return &RingBuffer{
		reader: r,
		data:   make([]byte, n),
	}
}

type RingBuffer struct {
	reader                  io.Reader
	dataLength              int
	readLength, writeLength int64
	rIdx, wIdx              int
	data                    []byte
}

func (rb *RingBuffer) Seek(offset int64, whence int) (actual int64, err error) {
	seeker, ok := rb.reader.(io.Seeker)

	if !ok {
		err = errors.Errorf("seeking not supported")
		return
	}

	if actual, err = seeker.Seek(offset, whence); err != nil {
		err = errors.Wrap(err)
		return
	}

	rb.dataLength = 0
	rb.readLength = actual
	rb.writeLength = actual
	rb.rIdx = 0
	rb.wIdx = 0

	if _, err = rb.Fill(); err != nil {
		if err == io.EOF {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	return
}

func (rb *RingBuffer) Reset(r io.Reader) {
	rb.reader = r
	rb.dataLength = 0
	rb.readLength = 0
	rb.writeLength = 0
	rb.rIdx = 0
	rb.wIdx = 0

	for i := range rb.data {
		rb.data[i] = 0
	}
}

func (rb *RingBuffer) ReadLength() int64 {
	return rb.readLength
}

func (rb *RingBuffer) PeekWriteable() (rs Slice) {
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

func (rb *RingBuffer) PeekReadable() (rs Slice) {
	if rb.Len() == 0 {
		return
	}

	rs.start = rb.readLength

	if rb.rIdx < rb.wIdx {
		rs.data[0] = rb.data[rb.rIdx:rb.wIdx]
	} else {
		rs.data[0] = rb.data[rb.rIdx:]
		rs.data[1] = rb.data[:rb.wIdx]
	}

	rCap := rs.Len()

	if rCap > rb.Len() && false {
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

func (rb *RingBuffer) PeekUnreadable() (rs Slice) {
	switch {
	case rb.wIdx < rb.rIdx:
		rs.data[0] = rb.data[rb.wIdx:rb.rIdx]

	case rb.rIdx < rb.wIdx:
		rs.data[0] = rb.data[:rb.rIdx]
		rs.data[1] = rb.data[rb.wIdx:]
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
		var f int64

		f, err = rb.Fill()

		switch {
		case err == io.EOF && f == 0:
			return

		case err != nil && err != io.EOF:
			return
		}
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

func (rb *RingBuffer) Fill() (n int64, err error) {
	if rb.reader == nil {
		panic("nil reader")
	}

	rs := rb.PeekWriteable()

	for i := 100; i > 0; i-- {
		n, err = rs.ReadFrom(rb.reader)
		rb.dataLength += int(n)
		rb.writeLength += n

		if int(n) <= rs.LenFirst() {
			rb.wIdx += int(n)
		} else {
			rb.wIdx = int(n) - len(rs.First())
		}

		if err != nil || n > 0 {
			return
		}
	}

	err = io.ErrNoProgress

	return
}

// func (rb *RingBuffer) SlideAndFill() (n int64, err error) {
// 	rb.data = append(rb.data, rb.PeekReadable().Second()...)
// 	copy(rb.data, rb.data[rb.rIdx:])
// 	rb.wIdx = rb.rIdx
// 	rb.rIdx = 0

// 	n, err = rb.Fill()

// 	return
// }

func (rb *RingBuffer) AdvanceRead(n int) {
	rb.rIdx += n
	rb.readLength += int64(n)

	if rb.rIdx > len(rb.data) {
		rb.rIdx -= len(rb.data)
	}

	rb.dataLength -= n
}

func (rb *RingBuffer) Unread(toUnread int) (actuallyUnread int) {
	if rb.wIdx < rb.rIdx {
		maxToUnread := rb.rIdx - rb.wIdx
		actuallyUnread = min(maxToUnread, toUnread)
		rb.rIdx -= actuallyUnread
	} else {
		actuallyUnread = min(toUnread, rb.rIdx)
		rb.rIdx -= actuallyUnread

		if actuallyUnread < toUnread && rb.rIdx == 0 {
			toUnread -= actuallyUnread
			last := rb.Cap() - 1
			toUnread = min(toUnread, last-rb.wIdx)
			rb.rIdx = last - toUnread
			actuallyUnread += toUnread
		}
	}

	rb.readLength -= int64(actuallyUnread)

	return
}

func (rb *RingBuffer) PeekUptoAndIncluding(b byte) (readable Slice, err error) {
	ok := false
	readable, ok = rb.PeekReadable().SliceUptoAndIncluding(b)

	if ok {
		return
	}

	_, err = rb.Fill()

	if errors.IsNotNilAndNotEOF(err) {
		return
	}

	err = nil
	readable, ok = rb.PeekReadable().SliceUptoAndIncluding(b)

	if !ok {
		err = collections.MakeErrNotFoundString(string(b))
		return
	}

	return
}

func (rb *RingBuffer) PeekUpto(b byte) (readable Slice, err error) {
	ok := false
	readable, ok = rb.PeekReadable().SliceUptoButExcluding(b)

	if ok {
		return
	}

	_, err = rb.Fill()

	if err != io.EOF && err != nil {
		return
	}

	readable, _ = rb.PeekReadable().SliceUptoButExcluding(b)

	return
}

func (rb *RingBuffer) Len() int {
	if rb.dataLength > len(rb.data) {
		panic("length is greater than buffer")
	}

	return rb.dataLength
}
