package catgut

import (
	"bytes"
	"io"
	"unicode"
	"unicode/utf8"
	"unsafe"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/ohio_buffer"
)

type String struct {
	addr *String
	Data bytes.Buffer
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

func (b *String) copyCheck() {
	if b.addr == nil {
		// This hack works around a failing of Go's escape analysis
		// that was causing b to escape and be heap allocated.
		// See issue 23382.
		// TODO-P1: once issue 7921 is fixed, this should be reverted to
		// just "b.addr = b".
		b.addr = (*String)(noescape(unsafe.Pointer(b)))
		return
	}

	if b.addr != b {
		oldB := b.addr
		b = &String{}
		b.addr = b
		b.Data = bytes.Buffer{}
		oldB.CopyTo(b)
		// panic("catgut: illegal use of non-zero String copied by value")
	}
}

func (str *String) String() string {
	str.copyCheck()
	return str.Data.String()
}

func (str *String) Len() int {
	str.copyCheck()
	return str.Data.Len()
}

func (str *String) AvailableBuffer() []byte {
	str.copyCheck()
	return str.Data.AvailableBuffer()
}

func (str *String) Bytes() []byte {
	str.copyCheck()
	return str.Data.Bytes()
}

func (str *String) Reset() {
	str.Data.Reset()
}

func (str *String) Write(p []byte) (int, error) {
	str.copyCheck()
	return str.Data.Write(p)
}

func (str *String) Grow(n int) {
	str.copyCheck()
	// c := str.Cap()
	str.Data.Grow(n)

	//	if c < str.Cap() {
	//		log.Debug().FunctionName(2)
	//	}
}

func (str *String) Map(mapping func(r rune) rune, s []byte) (n int) {
	str.Grow(len(s))
	b := str.AvailableBuffer()

	for i := 0; i < len(s); {
		wid := 1
		r := rune(s[i])

		if r >= utf8.RuneSelf {
			r, wid = utf8.DecodeRune(s[i:])
		}

		r = mapping(r)

		if r >= 0 {
			b = utf8.AppendRune(b, r)
		}

		i += wid
	}

	n, _ = str.Write(b)

	return
}

func (str *String) WriteLowerOrError(s []byte) (err error) {
	n := str.WriteLower(s)
	return ohio_buffer.MakeErrLength(int64(len(s)), int64(n), nil)
}

// WriteLower writes all Unicode letters mapped to
// their lower case.
func (str *String) WriteLower(s []byte) (n int) {
	isASCII, hasUpper := true, false

	for i := 0; i < len(s); i++ {
		c := s[i]

		if c >= utf8.RuneSelf {
			isASCII = false
			break
		}

		hasUpper = hasUpper || ('A' <= c && c <= 'Z')
	}

	str.Grow(len(s))

	if isASCII { // optimize for ASCII-only byte slices.
		b := str.AvailableBuffer()

		if !hasUpper {
			n, _ = str.Write(append(b, s...))
			return
		}

		for i := 0; i < len(s); i++ {
			c := s[i]

			if 'A' <= c && c <= 'Z' {
				c += 'a' - 'A'
			}

			b = append(b, c)
		}

		n, _ = str.Write(b)

		return
	}

	return str.Map(unicode.ToLower, s)
}

func (dst *String) Set(src string) (err error) {
	dst.copyCheck()
	dst.Reset()
	dst.Grow(len(src))

	b := append(dst.AvailableBuffer(), src...)
	dst.Write(b)

	return
}

func (src *String) CopyTo(dst *String) (err error) {
	dst.Reset()
	dst.Grow(src.Len())

	var n int

	n, err = dst.Write(src.Bytes())

	return ohio_buffer.MakeErrLength(int64(src.Len()), int64(n), err)
}

func (src *String) WriteTo(w io.Writer) (n int64, err error) {
	var n1 int
	n1, err = w.Write(src.Bytes())
	n = int64(n1)
	return
}

func (src *String) WriteToStringWriter(
	w schnittstellen.WriterAndStringWriter,
) (n int64, err error) {
	return src.WriteTo(w)
}