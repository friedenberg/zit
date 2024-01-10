package ennui_shas

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	ennui_upstream "github.com/friedenberg/zit/src/golf/ennui"
)

type Loc struct {
	Sha
	ennui_upstream.Range
}

func (l *Loc) Size() int {
	return l.Sha.Size() + l.Range.Size()
}

func (l *Loc) IsEmpty() bool {
	return l.Sha.IsNull() && l.Offset == 0 && l.ContentLength == 0
}

func (l *Loc) String() string {
	return fmt.Sprintf("%s@%s", &l.Sha, l.Range)
}

func (l *Loc) ReadFrom(r io.Reader) (n int64, err error) {
	var n2 int64
	n2, err = l.Sha.ReadFrom(r)
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = l.Range.ReadFrom(r)
	n += n2

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	return
}

func (l *Loc) WriteTo(w io.Writer) (n int64, err error) {
	var n2 int64
	n2, err = l.Sha.WriteTo(w)
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = l.Range.WriteTo(w)
	n += n2

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	return
}
