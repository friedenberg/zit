package etiketten_path

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/ohio"
)

//go:generate stringer -type=Type
type Type byte

// describe these
const (
	TypeDirect = Type(iota)
	TypeSuper
	TypeIndirect
)

func (t Type) ReadByte() (byte, error) {
	return byte(t), nil
}

func (t *Type) ReadFrom(r io.Reader) (n int64, err error) {
	var b [1]byte

	var n1 int
	n1, err = ohio.ReadAllOrDieTrying(r, b[:])
	n = int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	*t = Type(b[0])

	return
}

func (t Type) WriteTo(w io.Writer) (n int64, err error) {
	var b byte

	if b, err = t.ReadByte(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var n1 int
	n1, err = ohio.WriteAllOrDieTrying(w, []byte{b})
	n = int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
