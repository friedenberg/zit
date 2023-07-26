package collections

import (
	"sort"
	"strings"

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

// func ContainsKey(
// 	id schnittstellen.Stringer,
// 	cs ...schnittstellen.ContainsKeyer,
// ) (ok bool) {
// 	for _, c := range cs {
// 		if c.ContainsKey(id.String()) {
// 			return true
// 		}
// 	}

// 	return false
// }

func Len(cs ...schnittstellen.Lenner) (n int) {
	for _, c := range cs {
		n += c.Len()
	}

	return
}

func AddClone[E any, EPtr interface {
	*E
	ResetWithPtr(*E)
}](
	c schnittstellen.Adder[EPtr],
) schnittstellen.FuncIter[EPtr] {
	return func(e EPtr) (err error) {
		var e1 E
		EPtr(&e1).ResetWithPtr((*E)(e))
		c.Add(&e1)
		return
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

func Map[E schnittstellen.Value[E], F schnittstellen.Value[F]](
	in schnittstellen.SetLike[E],
	tr schnittstellen.FuncTransform[E, F],
) (out schnittstellen.MutableSetLike[F], err error) {
	out = MakeMutableSetStringer[F]()

	if err = in.Each(
		func(e E) (err error) {
			var e1 F

			if e1, err = tr(e); err != nil {
				if IsStopIteration(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
				}

				return
			}

			return out.Add(e1)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func DerivedValues[E schnittstellen.Value[E], F any](
	c schnittstellen.SetLike[E],
	f schnittstellen.FuncTransform[E, F],
) (out []F, err error) {
	out = make([]F, 0, c.Len())

	if err = c.Each(
		func(e E) (err error) {
			var e1 F

			if e1, err = f(e); err != nil {
				if IsStopIteration(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
				}

				return
			}

			out = append(out, e1)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func SortedValuesBy[E schnittstellen.Value[E]](
	c schnittstellen.SetLike[E],
	sf func(E, E) bool,
) (out []E) {
	out = c.Elements()

	sort.Slice(out, func(i, j int) bool { return sf(out[i], out[j]) })

	return
}

func SortedValues[E schnittstellen.Value[E]](
	c schnittstellen.SetLike[E],
) (out []E) {
	out = c.Elements()

	sort.Slice(out, func(i, j int) bool { return out[i].String() < out[j].String() })

	return
}

func Strings[E schnittstellen.ValueLike](
	c schnittstellen.SetLike[E],
) (out []string) {
	out = make([]string, 0, c.Len())

	c.Each(
		func(e E) (err error) {
			out = append(out, e.String())
			return
		},
	)

	return
}

func SortedStrings[E schnittstellen.ValueLike](
	c schnittstellen.SetLike[E],
) (out []string) {
	out = Strings(c)

	sort.Strings(out)

	return
}

func StringDelimiterSeparated[E schnittstellen.Value[E]](
	c schnittstellen.SetLike[E],
	d string,
) string {
	if c == nil {
		return ""
	}

	sorted := SortedStrings[E](c)

	sb := &strings.Builder{}
	first := true

	for _, e1 := range sorted {
		if !first {
			sb.WriteString(d)
		}

		sb.WriteString(e1)

		first = false
	}

	return sb.String()
}

func StringCommaSeparated[E schnittstellen.Value[E]](
	c schnittstellen.SetLike[E],
) string {
	return StringDelimiterSeparated(c, ", ")
}

func ReverseSortable(s sort.Interface) {
	max := s.Len() / 2

	for i := 0; i < max; i++ {
		s.Swap(i, s.Len()-1-i)
	}
}

func MakeFuncTransformer[T any, T1 any](
	wf schnittstellen.FuncIter[T],
) schnittstellen.FuncIter[T1] {
	return func(e T1) (err error) {
		if e1, ok := any(e).(T); ok {
			return wf(e1)
		}

		return
	}
}
