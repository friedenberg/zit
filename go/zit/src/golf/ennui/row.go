package ennui

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/sha"
)

const RowSize = sha.ByteSize + 1 + binary.MaxVarintLen64 + binary.MaxVarintLen64

type int64Bytes [binary.MaxVarintLen64]byte

type row struct {
	sha sha.Sha
	Loc
}

func (r *row) IsEmpty() bool {
	return r.Loc.IsEmpty() && r.sha.IsNull()
}

func (r *row) String() string {
	return fmt.Sprintf(
		"%s %s",
		&r.Loc,
		r.sha.GetShaString()[:4],
	)
}

func (current *row) ReadFrom(r io.Reader) (n int64, err error) {
	var n1 int64

	n1, err = current.sha.ReadFrom(r)
	n += n1

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	n1, err = current.Loc.ReadFrom(r)
	n += int64(n1)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	if n != RowSize {
		err = errors.Errorf("expected to read %d but read %d", RowSize, n)
		return
	}

	return
}

func (r *row) WriteTo(w io.Writer) (n int64, err error) {
	var n1 int
	var n2 int64

	n, err = r.sha.WriteTo(w)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = r.Loc.WriteTo(w)
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if n != RowSize {
		err = errors.Errorf("expected to write %d but read %d", RowSize, n)
		return
	}

	return
}

type rowEqualerComplete struct{}

func (rowEqualerComplete) Equals(a, b *row) bool {
	return a.sha.Equals(&b.sha) &&
		a.Loc.Page == b.Loc.Page &&
		a.Loc.Offset == b.Loc.Offset &&
		a.Loc.ContentLength == b.Loc.ContentLength
}

type rowEqualerShaOnly struct{}

func (rowEqualerShaOnly) Equals(a, b *row) bool {
	return a.sha.Equals(&b.sha)
}

type rowResetter struct{}

func (rowResetter) Reset(a *row) {
	a.sha.Reset()
	a.Page = 0
	a.Offset = 0
	a.ContentLength = 0
}

func (rowResetter) ResetWith(a, b *row) {
	a.sha.ResetWith(&b.sha)
	a.Page = b.Page
	a.Offset = b.Offset
	a.ContentLength = b.ContentLength
}

type rowLessor struct{}

func (rowLessor) Less(a, b *row) bool {
	cmp := bytes.Compare(a.sha.GetShaBytes(), b.sha.GetShaBytes())

	if cmp != 0 {
		return cmp == -1
	}

	if a.Page != b.Page {
		return a.Page < b.Page
	}

	if a.Offset != b.Offset {
		return a.Offset < b.Offset
	}

	return a.ContentLength < b.ContentLength
}
