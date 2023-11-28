package ohio_ring_buffer2

import (
	"bytes"
	"io"
	"strings"
)

type Slice struct {
	start int64
	data  [2][]byte
}

func (rs Slice) ReadFrom(r io.Reader) (n int64, err error) {
	var n1 int

	n1, err = r.Read(rs.First())
	n += int64(n1)

	if err != nil {
		return
	}

	n1, err = r.Read(rs.Second())
	n += int64(n1)

	if err != nil {
		return
	}

	return
}

func (rs Slice) WriteTo(w io.Writer) (n int64, err error) {
	var n1 int

	n1, err = w.Write(rs.First())
	n += int64(n1)

	if err != nil {
		return
	}

	n1, err = w.Write(rs.Second())
	n += int64(n1)

	return
}

func (rs Slice) Bytes() []byte {
	switch {
	case len(rs.First()) == rs.Len():
		return rs.First()

	case len(rs.Second()) == rs.Len():
		return rs.Second()

	default:
		var b bytes.Buffer
		b.Grow(rs.Len())
		b.Write(rs.First())
		b.Write(rs.Second())
		return b.Bytes()
	}
}

func (rs Slice) String() string {
	var s strings.Builder
	s.Grow(rs.Len())
	s.Write(rs.First())
	s.Write(rs.Second())
	return s.String()
}

func (rs Slice) Equal(b []byte) bool {
	c := 0

	for _, v := range rs.First() {
		if b[c] != v {
			return false
		}

		c++
	}

	for _, v := range rs.Second() {
		if b[c] != v {
			return false
		}

		c++
	}

	return true
}

func (rs Slice) Compare(b []byte) (c int) {
	if len(b) < len(rs.First()) {
		c = bytes.Compare(b, rs.First())
	} else {
		c = bytes.Compare(b[:len(rs.First())], rs.First())
	}

	if len(b) > len(rs.First()) && c == 0 {
		c = bytes.Compare(b[len(rs.First()):], rs.Second())
	}

	return
}

func (rs Slice) Start() int64 {
	return rs.start
}

func (rs Slice) First() []byte {
	return rs.data[0]
}

func (rs Slice) Second() []byte {
	return rs.data[1]
}

func (rs Slice) IsEmpty() bool {
	return rs.Len() == 0
}

func (rs Slice) Len() int {
	return len(rs.First()) + len(rs.Second())
}

func (rs Slice) Cut(b byte) (before, after Slice, ok bool) {
	for i, v := range rs.First() {
		if v == b {
			before = Slice{
				start: rs.start,
				data: [2][]byte{
					rs.First()[:i],
					nil,
				},
			}

			after = Slice{
				start: rs.start + int64(1+before.Len()),
				data: [2][]byte{
					rs.First()[i+1:],
					rs.Second(),
				},
			}

			ok = true

			return
		}
	}

	for i, v := range rs.Second() {
		if v == b {
			before = Slice{
				start: rs.start,
				data: [2][]byte{
					rs.First(),
					rs.Second()[:i],
				},
			}

			after = Slice{
				start: rs.start + int64(1+before.Len()),
				data: [2][]byte{
					rs.Second()[i+1:],
					nil,
				},
			}

			ok = true

			return
		}
	}

	return
}

func (rs Slice) Upto(b byte) (s Slice, ok bool) {
	for i, v := range rs.First() {
		if v == b {
			s = Slice{
				start: rs.start,
				data: [2][]byte{
					rs.First()[:i],
					nil,
				},
			}

			ok = true

			return
		}
	}

	for i, v := range rs.Second() {
		if v == b {
			s = Slice{
				start: rs.start,
				data: [2][]byte{
					rs.First(),
					rs.Second()[:i],
				},
			}

			ok = true

			return
		}
	}

	return
}

func (rs Slice) findFromStart(ff FindFunc) (length int, partial bool) {
	var offset int

	offset, length, partial = ff(rs)

	if offset > 0 {
		length = 0
		partial = false
	}

	return
}

func (rs Slice) findAnywhere(ff FindFunc) (offset, length int, partial bool) {
	if rs.Len() == 0 {
		return
	}

	return ff(rs)
}
