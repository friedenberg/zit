package kennung

import (
	"regexp"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type expanderAll struct {
	delimiter *regexp.Regexp
}

func MakeExpanderAll(
	delimiter string,
) expanderAll {
	return expanderAll{
		delimiter: regexp.MustCompile(delimiter),
	}
}

func (ex expanderAll) Expand(
	sa schnittstellen.FuncSetString,
	s string,
) {
	sa(s)

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

		sa(t1)
		sa(t2)

		if 0 < i && i < len(hyphens) {
			t1 := s[prevLocEnd:locStart]
			sa(t1)
		}

		prevLocEnd = locEnd
	}

	return
}