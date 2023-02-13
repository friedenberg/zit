package gattung

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

var ErrNoAbbreviation = errors.New("no abbreviation")

func MakeErrUnsupportedGattung(g schnittstellen.GattungGetter) error {
	return errors.WrapN(1, errUnsupportedGattung{Gattung: g.GetGattung()})
}

func IsErrUnsupportedGattung(err error) bool {
	return errors.Is(err, errUnsupportedGattung{Gattung: Unknown})
}

type errUnsupportedGattung struct {
	schnittstellen.Gattung
}

func (e errUnsupportedGattung) Is(target error) (ok bool) {
	_, ok = target.(errUnsupportedGattung)
	return
}

func (e errUnsupportedGattung) Error() string {
	return fmt.Sprintf("unsupported gattung: %q", e.Gattung)
}

type ErrUnrecognizedGattung struct {
	string
}

func (e ErrUnrecognizedGattung) Is(target error) (ok bool) {
	_, ok = target.(ErrUnrecognizedGattung)
	return
}

func (e ErrUnrecognizedGattung) Error() string {
	return fmt.Sprintf("unknown gattung: %q", e.string)
}

type ErrWrongType struct {
	ExpectedType, ActualType Gattung
}

func (e ErrWrongType) Is(target error) (ok bool) {
	_, ok = target.(ErrWrongType)
	return
}

func (e ErrWrongType) Error() string {
	return fmt.Sprintf("expected zk_types %s but got %s", e.ExpectedType, e.ActualType)
}

type ErrEmptyKennung struct{}

func (e ErrEmptyKennung) Is(target error) (ok bool) {
	_, ok = target.(ErrEmptyKennung)
	return
}

func (e ErrEmptyKennung) Error() string {
	return fmt.Sprintf("empty kennung")
}
