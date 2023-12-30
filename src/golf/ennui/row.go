package ennui

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/ohio"
)

const RowSize = sha.ByteSize + 1 + binary.MaxVarintLen64

type int64Bytes [binary.MaxVarintLen64]byte

type row struct {
	sha    sha.Sha
	page   [1]byte
	offset int64Bytes
}

func (r *row) getOffset() int64 {
	o, _ := binary.Varint(r.offset[:])
	return o
}

func (r *row) getPage() uint8 {
	return uint8(r.page[0])
}

func (r *row) setOffset(v int64) {
	binary.PutVarint(r.offset[:], v)

	if v != r.getOffset() {
		panic(fmt.Sprintf("expected %d but got %d", v, r.getOffset()))
	}
}

func (r *row) String() string {
	return fmt.Sprintf("%d@%d -> %s", uint8(r.page[0]), r.getOffset(), &r.sha)
}

func (current *row) ReadFrom(r io.Reader) (n int64, err error) {
	var n1 int64
	var n2 int

	n1, err = current.sha.ReadFrom(r)
	n += n1

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	n2, err = ohio.ReadAllOrDieTrying(r, current.page[:])
	n += int64(n2)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	if current.getPage() > 16 {
		err = errors.Errorf("page too big: %s", current)
		return
	}

	n2, err = ohio.ReadAllOrDieTrying(r, current.offset[:])
	n += int64(n2)

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
	var n2 int

	n, err = r.sha.WriteTo(w)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = ohio.WriteAllOrDieTrying(w, r.page[:])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = ohio.WriteAllOrDieTrying(w, r.offset[:])
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

type rowEqualer struct{}

func (rowEqualer) Equals(a, b *row) bool {
	return a.sha.Equals(&b.sha) &&
		bytes.Equal(a.page[:], b.page[:]) &&
		bytes.Equal(a.offset[:], b.offset[:])
}

type rowResetter struct{}

func (rowResetter) Reset(a *row) {
	a.sha.Reset()
	a.page[0] = 0
	a.offset = [binary.MaxVarintLen64]byte{}
}

func (rowResetter) ResetWith(a, b *row) {
	a.sha.ResetWith(&b.sha)
	a.page[0] = b.page[0]
	a.offset = b.offset
}

type rowLessor struct{}

func (rowLessor) Less(a, b *row) bool {
	cmp := bytes.Compare(a.sha.GetShaBytes(), b.sha.GetShaBytes())

	if cmp == 0 {
		cmp = bytes.Compare(a.page[:], b.page[:])
	}

	if cmp == 0 {
		cmp = bytes.Compare(a.offset[:], b.offset[:])
	}

	return cmp == -1
}
