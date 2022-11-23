package string_expansion

import (
	"regexp"

	"github.com/friedenberg/zit/src/bravo/collections"
)

type expanderAll[T collections.ValueElement, T1 collections.ValueElementPtr[T]] struct {
	delimiter *regexp.Regexp
}

func MakeExpanderAll[T collections.ValueElement, T1 collections.ValueElementPtr[T]](
	delimiter string,
) expanderAll[T, T1] {
	return expanderAll[T, T1]{
		delimiter: regexp.MustCompile(delimiter),
	}
}

func (ex expanderAll[T, T1]) Expand(s string) (out collections.ValueSet[T, T1]) {
	expanded := collections.MakeMutableValueSet[T, T1]()
	expanded.AddString(s)

	defer func() {
		out = expanded.Copy()
	}()

	if s == "" {
		return
	}

	hyphens := ex.delimiter.FindAllIndex([]byte(s), -1)

	if hyphens == nil {
		return
	}

	end := len(s)
	prevLocEnd := 0

	for i, loc := range hyphens {
		locStart := loc[0]
		locEnd := loc[1]
		t1 := s[0:locStart]
		t2 := s[locEnd:end]

		expanded.AddString(t1)
		expanded.AddString(t2)

		if 0 < i && i < len(hyphens) {
			t1 := s[prevLocEnd:locStart]
			expanded.AddString(t1)
		}

		prevLocEnd = locEnd
	}

	return
}
