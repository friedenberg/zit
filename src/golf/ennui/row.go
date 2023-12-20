package ennui

import (
	"bytes"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type row [2]sha.Sha

func (r *row) String() string {
	return fmt.Sprintf("%s %s", &r[0], &r[1])
}

func (current *row) ReadFrom(r io.Reader) (n int64, err error) {
	var n1 int64
	n1, err = current[0].ReadFrom(r)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = current[1].ReadFrom(r)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (r *row) WriteTo(w io.Writer) (n int64, err error) {
	var n1 int

	n, err = r[0].WriteTo(w)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n, err = r[1].WriteTo(w)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type rowEqualer struct{}

func (rowEqualer) Equals(a, b *row) bool {
	zeroEquals := a[0].Equals(&b[0])
	oneEquals := a[1].Equals(&b[1])
	// log.Debug().Caller(1, "%s =? %s: %t, %t", a, b, zeroEquals, oneEquals)
	return zeroEquals && oneEquals
}

type rowResetter struct{}

func (rowResetter) Reset(a *row) {
	a[0].Reset()
	a[1].Reset()
}

func (rowResetter) ResetWith(a, b *row) {
	a[0].ResetWith(&b[0])
	a[1].ResetWith(&b[1])
}

type rowLessor struct{}

func (rowLessor) Less(a, b *row) bool {
	cmp := bytes.Compare(a[0].GetShaBytes(), b[0].GetShaBytes())

	if cmp == 0 {
		cmp = bytes.Compare(a[1].GetShaBytes(), b[1].GetShaBytes())
	}

	return cmp == -1
}
