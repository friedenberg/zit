package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type nopAkteFormat[
	T Objekte[T],
	T1 ObjektePtr[T],
] struct{}

func MakeNopAkteFormat[
	T Objekte[T],
	T1 ObjektePtr[T],
]() nopAkteFormat[T, T1] {
	return nopAkteFormat[T, T1]{}
}

func (_ nopAkteFormat[T, T1]) Parse(r io.Reader, _ T1) (n int64, err error) {
	if n, err = io.Copy(io.Discard, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (_ nopAkteFormat[T, T1]) Format(w io.Writer, _ T1) (n int64, err error) {
	errors.TodoP0("how to format without content?")
	return
}
