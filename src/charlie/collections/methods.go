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
	c schnittstellen.Adder[E],
) schnittstellen.FuncSetString {
	return func(v string) (err error) {
		return AddString[E, EPtr](c, v)
	}
}

func AddString[E any, EPtr schnittstellen.SetterPtr[E]](
	c schnittstellen.Adder[E],
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

func ExpandAndAddString[E any, EPtr schnittstellen.SetterPtr[E]](
	c schnittstellen.Adder[E],
	expander func(string) (string, error),
	v string,
) (err error) {
	if expander != nil {
		v1 := v

		if v1, err = expander(v); err != nil {
			err = nil
			v1 = v
		}

		v = v1
	}

	if err = AddString[E, EPtr](c, v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type AddGetKeyer[E schnittstellen.Lessor[E]] interface {
	schnittstellen.Adder[E]
	Get(string) (E, bool)
	Key(E) string
}

func AddIfGreater[E schnittstellen.Lessor[E]](
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
	c schnittstellen.EachPtrer[E],
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
