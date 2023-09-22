package iter

import (
	"sort"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

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

func SortedValuesBy[E any](
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

	sort.Slice(
		out,
		func(i, j int) bool { return out[i].String() < out[j].String() },
	)

	return
}

func Strings[E schnittstellen.Stringer](
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

func SortedStrings[E schnittstellen.Stringer](
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
