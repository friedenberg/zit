package store_verzeichnisse

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/delta/schlussel"
)

type binaryField struct {
	schlussel.Schlussel
	ContentLength [2]uint8
	Content       bytes.Buffer
}

func (bf *binaryField) String() string {
	cl, _, _ := bf.GetContentLength()
	return fmt.Sprintf("%s:%d:%x", bf.Schlussel, cl, bf.Content.Bytes())
}

func (bf *binaryField) Reset() {
	bf.Schlussel.Reset()
	bf.ContentLength[0] = 0
	bf.ContentLength[1] = 0
	bf.Content.Reset()
}

func (bf *binaryField) GetContentLength() (contentLength int, contentLength64 int64, err error) {
	var n int
	contentLength64, n = binary.Varint(bf.ContentLength[:])

	if n <= 0 {
		err = errors.Errorf("error in content length: %d", n)
		return
	}

	if contentLength64 > math.MaxUint16 {
		err = errContentLengthTooLarge
		return
	}

	if contentLength64 < 0 {
		err = errContentLengthNegative
		return
	}

	return int(contentLength64), contentLength64, nil
}

func (bf *binaryField) SetContentLength(v int) {
	if v < 0 {
		panic(errContentLengthNegative)
	}

	if v > math.MaxUint16 {
		panic(errContentLengthTooLarge)
	}

	// TODO
	binary.PutVarint(bf.ContentLength[:], int64(v))
}

var (
	errContentLengthTooLarge = errors.New("content length too large")
	errContentLengthNegative = errors.New("content length negative")
)

func (bf *binaryField) ReadFrom(r io.Reader) (n int64, err error) {
	var n1 int
	var n2 int64
	n2, err = bf.Schlussel.ReadFrom(r)
	n += int64(n2)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	n1, err = ohio.ReadAllOrDieTrying(r, bf.ContentLength[:])
	n += int64(n1)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	contentLength, contentLength64, err := bf.GetContentLength()
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	bf.Content.Grow(contentLength)
	bf.Content.Reset()

	n2, err = io.CopyN(&bf.Content, r, contentLength64)
	n += n2

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	return
}

var errContentLengthDoesNotMatchContent = errors.New(
	"content length does not match content",
)

func (bf *binaryField) WriteTo(w io.Writer) (n int64, err error) {
	if bf.Content.Len() > math.MaxUint16 {
		err = errContentLengthTooLarge
		return
	}

	bf.SetContentLength(bf.Content.Len())

	var n1 int
	var n2 int64
	n2, err = bf.Schlussel.WriteTo(w)
	n += int64(n2)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	n1, err = ohio.WriteAllOrDieTrying(w, bf.ContentLength[:])
	n += int64(n1)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	n2, err = io.Copy(w, &bf.Content)
	n += n2

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	return
}
