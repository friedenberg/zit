package iter

import (
	"sort"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
)

func SortedValuesBy[E any](
	c schnittstellen.SetLike[E],
	sf func(E, E) bool,
) (out []E) {
	out = Elements(c)

	sort.Slice(out, func(i, j int) bool { return sf(out[i], out[j]) })

	return
}

func SortedValues[E schnittstellen.Value[E]](
	c schnittstellen.SetLike[E],
) (out []E) {
	out = Elements(c)

	sort.Slice(
		out,
		func(i, j int) bool { return out[i].String() < out[j].String() },
	)

	return
}

func Strings[E schnittstellen.Stringer](
	cs ...schnittstellen.SetLike[E],
) (out []string) {
	l := 0

	for _, c := range cs {
		if c == nil {
			continue
		}

		l += c.Len()
	}

	out = make([]string, 0, l)

	for _, c := range cs {
		if c == nil {
			continue
		}

		err := c.Each(
			func(e E) (err error) {
				out = append(out, e.String())
				return
			},
		)

		errors.PanicIfError(err)
	}

	return
}

func SortedStrings[E schnittstellen.Stringer](
	cs ...schnittstellen.SetLike[E],
) (out []string) {
	out = Strings(cs...)

	sort.Strings(out)

	return
}

func StringDelimiterSeparated[E schnittstellen.Stringer](
	d string,
	cs ...schnittstellen.SetLike[E],
) string {
	if cs == nil {
		return ""
	}

	sorted := SortedStrings[E](cs...)

	if len(sorted) == 0 {
		return ""
	}

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

func StringCommaSeparated[E schnittstellen.Stringer](
	cs ...schnittstellen.SetLike[E],
) string {
	return StringDelimiterSeparated(", ", cs...)
}

func ReverseSortable(s sort.Interface) {
	max := s.Len() / 2

	for i := 0; i < max; i++ {
		s.Swap(i, s.Len()-1-i)
	}
}
