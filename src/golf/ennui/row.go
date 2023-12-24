package ennui

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
)

const RowSize = sha.ByteSize + 1 + binary.MaxVarintLen64

type int64Bytes [binary.MaxVarintLen64]byte

type row struct {
	sha    sha.Sha
	page   [1]byte
	offset int64Bytes
}

func (r *row) String() string {
	return fmt.Sprintf("%s %d", &r.sha, r.page)
}

func (current *row) ReadFrom(r io.Reader) (n int64, err error) {
	var n1 int64
	var n2 int

	n1, err = current.sha.ReadFrom(r)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = r.Read(current.page[:])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = r.Read(current.offset[:])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
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

	n2, err = w.Write(r.page[:])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = w.Write(r.offset[:])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
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
