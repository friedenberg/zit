package object_probe_index

import (
	"encoding/binary"
	"fmt"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
)

type Range struct {
	Offset, ContentLength int64
}

func (l *Range) Size() int {
	return binary.MaxVarintLen64 * 2
}

func (l Range) IsEmpty() bool {
	return l.Offset == 0 && l.ContentLength == 0
}

func (l Range) String() string {
	return fmt.Sprintf("%03d+%03d", l.Offset, l.ContentLength)
}

func (l *Range) ReadFrom(r io.Reader) (n int64, err error) {
	var n1 int
	l.Offset, n1, err = ohio.ReadInt64(r)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	l.ContentLength, n1, err = ohio.ReadInt64(r)
	n += int64(n1)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	return
}

func (l *Range) WriteTo(w io.Writer) (n int64, err error) {
	var n1 int
	n1, err = ohio.WriteInt64(w, l.Offset)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = ohio.WriteInt64(w, l.ContentLength)
	n += int64(n1)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	return
}
