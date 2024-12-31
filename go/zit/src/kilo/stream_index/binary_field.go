package stream_index

import (
	"bytes"
	"fmt"
	"io"
	"math"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/delta/keys"
)

type binaryField struct {
	keys.Binary
	ContentLength uint16
	Content       bytes.Buffer
}

func (bf *binaryField) String() string {
	return fmt.Sprintf("%s:%d:%x", bf.Binary, bf.ContentLength, bf.Content.Bytes())
}

func (bf *binaryField) Reset() {
	bf.Binary.Reset()
	bf.ContentLength = 0
	bf.Content.Reset()
}

var (
	errContentLengthTooLarge = errors.New("content length too large")
	errContentLengthNegative = errors.New("content length negative")
)

func (bf *binaryField) ReadFrom(r io.Reader) (n int64, err error) {
	var n1 int
	var n2 int64
	n2, err = bf.Binary.ReadFrom(r)
	n += int64(n2)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	n1, bf.ContentLength, err = ohio.ReadFixedUInt16(r)
	n += int64(n1)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	bf.Content.Grow(int(bf.ContentLength))
	bf.Content.Reset()

	n2, err = io.CopyN(&bf.Content, r, int64(bf.ContentLength))
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
	// defer func() {
	// 	r := recover()

	// 	if r == nil {
	// 		return
	// 	}

	// 	ui.Debug().Print(bf.Content.Len(), bf.ContentLength)
	// 	panic(r)
	// }()

	if bf.Content.Len() > math.MaxUint16 {
		err = errContentLengthTooLarge
		return
	}

	bf.ContentLength = uint16(bf.Content.Len())

	var n1 int
	var n2 int64
	n2, err = bf.Binary.WriteTo(w)
	n += int64(n2)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	n1, err = ohio.WriteFixedUInt16(w, bf.ContentLength)
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
