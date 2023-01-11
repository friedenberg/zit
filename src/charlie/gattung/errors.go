package gattung

import (
	"errors"
	"fmt"
)

var (
	ErrUnsupportedGattung = errors.New("unsupported gattung")
)

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

type ErrEmptyKennung struct {
}

func (e ErrEmptyKennung) Is(target error) (ok bool) {
	_, ok = target.(ErrEmptyKennung)
	return
}

func (e ErrEmptyKennung) Error() string {
	return fmt.Sprintf("empty kennung")
}
