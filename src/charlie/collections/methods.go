package collections

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func MakeFuncSetString[
	E any,
	EPtr schnittstellen.SetterPtr[E],
](
	c Adder[E],
) schnittstellen.FuncSetString {
	return func(v string) (err error) {
		return AddString[E, EPtr](c, v)
	}
}

func AddString[E any, EPtr schnittstellen.SetterPtr[E]](
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

type AddGetKeyer[E Lessor[E]] interface {
	Adder[E]
	Get(string) (E, bool)
	Key(E) string
}

func AddIfGreater[E Lessor[E]](
	c AddGetKeyer[E],
	e E,
) (ok bool) {
	k := c.Key(e)
	var old E

	if old, ok = c.Get(k); !ok || old.Less(e) {
		c.Add(e)
	}

	return
}

func String[E schnittstellen.Value](
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
