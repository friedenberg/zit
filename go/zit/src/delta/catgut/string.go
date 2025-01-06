package catgut

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"
	"unsafe"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio_buffer"
)

type String struct {
	addr *String
	data bytes.Buffer
}

func MakeFromReader(r io.Reader, limit int) (s *String, err error) {
	s = GetPool().Get()

	if _, err = s.ReadNFrom(r, limit); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeFromString(v string) (s *String) {
	s = GetPool().Get()
	errors.PanicIfError(s.Set(v))
	return
}

func MakeFromBytes(b []byte) (s *String) {
	s = GetPool().Get()
	errors.PanicIfError(s.SetBytes(b))
	return s
}

func Make(b *String) (a *String) {
	a = GetPool().Get()
	errors.PanicIfError(a.SetBytes(b.Bytes()))
	return
}

// noescape hides a pointer from escape analysis. It is the identity function
// but escape analysis doesn't think the output depends on the input.
// noescape is inlined and currently compiles down to zero instructions.
// USE CAREFULLY!
// This was copied from the runtime; see issues 23382 and 7921.
//
//go:nosplit
//go:nocheckptr
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

const (
	configPanicOnCopy    = true
	configDebugCopyCheck = false
)

func (b *String) copyCheck() {
	if b.addr == nil {
		// This hack works around a failing of Go's escape analysis
		// that was causing b to escape and be heap allocated.
		// See issue 23382.
		// TODO-P1: once issue 7921 is fixed, this should be reverted to
		// just "b.addr = b".
		if configDebugCopyCheck {
			ui.Log().Caller(6, "saved addr: %d", unsafe.Pointer(b))
		}
		b.addr = (*String)(noescape(unsafe.Pointer(b)))
		return
	}

	if b.addr != b {
		if configPanicOnCopy {
			panic(
				fmt.Sprintf(
					"catgut: illegal use of non-zero String copied by value: %d",
					unsafe.Pointer(b.addr),
				),
			)
		} else {
			oldB := b.addr
			b = &String{}
			b.addr = b
			b.data = bytes.Buffer{}
			oldB.CopyTo(b)
		}
	}
}

func (a *String) ComparePartialOffset(b *String, offset int) int {
	lenA, lenB := a.Len(), b.Len()

	if offset > min(lenA, lenB)-1 {
		panic("offset out of bounds")
	}

	return CompareUTF8Bytes(
		ComparerBytes(a.Bytes()[offset:]),
		ComparerBytes(b.Bytes()[offset:]),
		true,
	)
}

func (a *String) ComparePartial(b *String) int {
	return CompareUTF8Bytes(
		ComparerBytes(a.Bytes()),
		ComparerBytes(b.Bytes()),
		true,
	)
}

func (a *String) ComparePartialComparer(b Comparer) int {
	return CompareUTF8Bytes(
		ComparerBytes(a.Bytes()),
		b,
		true,
	)
}

func (a *String) Compare(b *String) int {
	return CompareUTF8Bytes(
		ComparerBytes(a.Bytes()),
		ComparerBytes(b.Bytes()),
		false,
	)
}

func (str *String) String() string {
	str.copyCheck()
	return str.data.String()
}

func (str *String) Len() int {
	str.copyCheck()
	return str.data.Len()
}

func (str *String) IsEmpty() bool {
	if str == nil {
		return true
	}

	str.copyCheck()
	return str.data.Len() == 0
}

func (a *String) Equals(b *String) bool {
	return bytes.Equal(a.Bytes(), b.Bytes())
}

func (a *String) EqualsString(b string) bool {
	return string(a.Bytes()) == b
}

func (a *String) EqualsBytes(b []byte) bool {
	return bytes.Equal(a.Bytes(), b)
}

func (str *String) AvailableBuffer() []byte {
	str.copyCheck()
	return str.data.AvailableBuffer()
}

func (str *String) Available() int {
	str.copyCheck()
	return str.data.Available()
}

func (str *String) Bytes() []byte {
	if str == nil {
		return nil
	}

	str.copyCheck()
	return str.data.Bytes()
}

func (str *String) Reset() {
	str.data.Reset()
}

func (str *String) Write(p []byte) (int, error) {
	str.copyCheck()
	return str.data.Write(p)
}

func (str *String) Append(vs ...*String) (n int, err error) {
	var n1 int

	for _, v := range vs {
		n1, err = str.Write(v.Bytes())
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (str *String) Grow(n int) {
	str.copyCheck()
	str.data.Grow(n)
}

func (str *String) Map(mapping func(r rune) rune, s []byte) (n int) {
	return MapTo(&str.data, s, mapping)
}

func (str *String) WriteLowerOrError(s []byte) (err error) {
	n := str.WriteLower(s)
	return ohio_buffer.MakeErrLength(int64(len(s)), int64(n), nil)
}

// WriteLower writes all Unicode letters mapped to
// their lower case.
func (str *String) WriteLower(s []byte) (n int) {
	return WriteLower(&str.data, s)
}

func (dst *String) WriteRune(r rune) (int, error) {
	dst.Grow(utf8.RuneLen(r))
	b := dst.data.AvailableBuffer()
	b = utf8.AppendRune(b, r)
	return dst.data.Write(b)
}

func (dst *String) SetBytes(src []byte) (err error) {
	dst.copyCheck()
	dst.Reset()
	dst.Grow(len(src))

	b := append(dst.AvailableBuffer(), src...)
	var n int
	n, err = dst.Write(b)

	if n != len(src) {
		panic(fmt.Sprintf("tried to write %d but only wrote %d", len(src), n))
	}

	return
}

func (dst *String) Set(src string) (err error) {
	dst.copyCheck()
	dst.Reset()
	dst.Grow(len(src))

	b := append(dst.AvailableBuffer(), src...)
	dst.Write(b)

	return
}

func (dst *String) ReadFromBuffer(src *bytes.Buffer) (err error) {
	dst.Reset()
	dst.Grow(src.Len())
	var n int
	n, err = dst.Write(append(dst.AvailableBuffer(), src.Bytes()...))
	return ohio_buffer.MakeErrLength(int64(src.Len()), int64(n), err)
}

func (dst *String) ReadFrom(r io.Reader) (n int64, err error) {
	dst.Reset()
	return dst.data.ReadFrom(r)
}

func (dst *String) ReadNFrom(r io.Reader, toRead int) (read int, err error) {
	dst.Reset()
	dst.Grow(toRead)
	b := dst.AvailableBuffer()[:toRead]

	read, err = ohio.ReadAllOrDieTrying(r, b)
	if err != nil {
		if read == toRead && err == io.EOF {
			err = nil
		} else {
			err = errors.WrapExcept(err, io.EOF)
			return
		}
	}

	if _, err = dst.data.Write(b); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (src *String) Clone() (dst *String, err error) {
	dst = GetPool().Get()

	if err = src.CopyTo(dst); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (src *String) CopyTo(dst *String) (err error) {
	dst.Reset()
	dst.Grow(src.Len())

	var n int

	n, err = dst.Write(src.Bytes())

	return MakeErrLength(int64(src.Len()), int64(n), err)
}

func (src *String) WriteTo(w io.Writer) (n int64, err error) {
	var n1 int
	n1, err = w.Write(src.Bytes())
	n = int64(n1)
	return
}

func (src *String) WriteToStringWriter(
	w interfaces.WriterAndStringWriter,
) (n int64, err error) {
	return src.WriteTo(w)
}

func (src *String) MarshalText() ([]byte, error) {
	return src.Bytes(), nil
}

func (src *String) UnmarshalText(b []byte) error {
	if src == nil {
		src = &String{}
	}

	src.SetBytes(b)
	return nil
}

func (src *String) MarshalBinary() ([]byte, error) {
	return src.Bytes(), nil
}

func (src *String) UnmarshalBinary(b []byte) error {
	if src == nil {
		src = &String{}
	}

	src.SetBytes(b)
	return nil
}
