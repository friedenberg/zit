package collections

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

type FuncSetString = gattung.FuncSetString

type Setter interface {
	Set(string) error
}

type SetterPtr[T any] interface {
	gattung.ElementPtr[T]
	Setter
}

type Adder[E any] interface {
	Add(E) error
}

type StringAdder interface {
	AddString(string) error
}

func MakeFuncSetString[E gattung.Element, EPtr SetterPtr[E]](
	c Adder[E],
) FuncSetString {
	return func(v string) (err error) {
		return AddString[E, EPtr](c, v)
	}
}

func AddString[E gattung.Element, EPtr SetterPtr[E]](
	c Adder[E],
	v string,
) (err error) {
	var e E

	if err = EPtr(&e).Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.Add(e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func ReverseSortable(s sort.Interface) {
	max := s.Len() / 2

	for i := 0; i < max; i++ {
		s.Swap(i, s.Len()-1-i)
	}
}
