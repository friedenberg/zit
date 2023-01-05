package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

type nopAkteParser[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
] struct {
}

func MakeNopAkteParser[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
]() nopAkteParser[T, T1] {
	return nopAkteParser[T, T1]{}
}

func (_ nopAkteParser[T, T1]) Parse(r io.Reader, _ T1) (n int64, err error) {
	if n, err = io.Copy(io.Discard, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
