package ennui

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Loc struct {
	Page   uint8
	Offset uint64
}

func (l *Loc) String() string {
	return fmt.Sprintf("%d@%d", l.Page, l.Offset)
}

func (l *Loc) ReadFrom(r io.Reader) (n int64, err error) {
	var page [1]byte

	_, err = r.Read(page[:])

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	var offset int64Bytes

	_, err = r.Read(offset[:])

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	var n1 int
	l.Page = page[0]
	l.Offset, n1 = binary.Uvarint(offset[:])

	if n1 <= 0 {
		err = errors.Errorf("not a valid uint64")
		return
	}

	return
}
