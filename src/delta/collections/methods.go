package collections

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

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

func String[E ValueSetElement](
	c EachPtrer[E],
) string {
	errors.TodoP1("implement")
	return ""
}

func ReverseSortable(s sort.Interface) {
	max := s.Len() / 2

	for i := 0; i < max; i++ {
		s.Swap(i, s.Len()-1-i)
	}
}
