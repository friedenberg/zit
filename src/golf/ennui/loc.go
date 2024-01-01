package ennui

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/ohio"
)

type Loc struct {
	Page                  uint8
	Offset, ContentLength int64
}

func (l Loc) IsEmpty() bool {
	return l.Page == 0 && l.Offset == 0 && l.ContentLength == 0
}

func (l Loc) String() string {
	return fmt.Sprintf("%02d@%03d", l.Page, l.Offset)
}

func (l *Loc) ReadFrom(r io.Reader) (n int64, err error) {
	var n1 int
	var intErr int
	var page [1]byte

	n1, err = ohio.ReadAllOrDieTrying(r, page[:])
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	var pageInt int64
	pageInt, intErr = binary.Varint(page[:])

	if intErr <= 0 {
		err = errors.Errorf("page parse issue: %d", intErr)
		return
	}

	if pageInt > 16 {
		err = errors.Errorf("page too big: %d", pageInt)
		return
	}

	l.Page = uint8(pageInt)

	var b int64Bytes

	n1, err = ohio.ReadAllOrDieTrying(r, b[:])
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	l.Offset, intErr = binary.Varint(b[:])

	if intErr <= 0 {
		err = errors.Errorf("offset parse issue: %d", intErr)
		return
	}

	n1, err = ohio.ReadAllOrDieTrying(r, b[:])
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	l.ContentLength, intErr = binary.Varint(b[:])

	if intErr <= 0 {
		err = errors.Errorf("content length parse issue: %d", intErr)
		return
	}

	return
}

func (l *Loc) WriteTo(w io.Writer) (n int64, err error) {
	var n1 int
	var intErr int
	var page [1]byte

	intErr = binary.PutVarint(page[:], int64(l.Page))

	if intErr != 1 {
		err = errors.Errorf("expected to write %d but wrote %d", 1, intErr)
		return
	}

	n1, err = ohio.WriteAllOrDieTrying(w, page[:])
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	var b int64Bytes

	binary.PutVarint(b[:], l.Offset)

	n1, err = ohio.WriteAllOrDieTrying(w, b[:])
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	binary.PutVarint(b[:], l.ContentLength)

	n1, err = ohio.WriteAllOrDieTrying(w, b[:])
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
