package gattung

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

var ErrNoAbbreviation = errors.New("no abbreviation")

func MakeErrUnsupportedGattung(g interfaces.GenreGetter) error {
	return errors.WrapN(1, errUnsupportedGattung{Genre: g.GetGenre()})
}

func IsErrUnsupportedGattung(err error) bool {
	return errors.Is(err, errUnsupportedGattung{Genre: Unknown})
}

type errUnsupportedGattung struct {
	interfaces.Genre
}

func (e errUnsupportedGattung) Is(target error) (ok bool) {
	_, ok = target.(errUnsupportedGattung)
	return
}

func (e errUnsupportedGattung) Error() string {
	return fmt.Sprintf("unsupported gattung: %q", e.Genre)
}

func MakeErrUnrecognizedGattung(v string) errUnrecognizedGattung {
	return errUnrecognizedGattung(v)
}

func IsErrUnrecognizedGattung(err error) bool {
	return errors.Is(err, errUnrecognizedGattung(""))
}

type errUnrecognizedGattung string

func (e errUnrecognizedGattung) Is(target error) (ok bool) {
	_, ok = target.(errUnrecognizedGattung)
	return
}

func (e errUnrecognizedGattung) Error() string {
	return fmt.Sprintf("unknown gattung: %q", string(e))
}

type ErrWrongType struct {
	ExpectedType, ActualType Genre
}

func (e ErrWrongType) Is(target error) (ok bool) {
	_, ok = target.(ErrWrongType)
	return
}

func (e ErrWrongType) Error() string {
	return fmt.Sprintf(
		"expected zk_types %s but got %s",
		e.ExpectedType,
		e.ActualType,
	)
}

type ErrEmptyKennung struct{}

func (e ErrEmptyKennung) Is(target error) (ok bool) {
	_, ok = target.(ErrEmptyKennung)
	return
}

func (e ErrEmptyKennung) Error() string {
	return fmt.Sprintf("empty kennung")
}
