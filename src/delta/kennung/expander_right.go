package kennung

import (
	"regexp"
)

type expanderRight struct {
	delimiter *regexp.Regexp
}

func MakeExpanderRight(
	delimiter string,
) expanderRight {
	return expanderRight{
		delimiter: regexp.MustCompile(delimiter),
	}
}

func (ex expanderRight) Expand(
	sa stringAdder,
	s string,
) {
	sa.AddString(s)

	if s == "" {
		return
	}

	hyphens := ex.delimiter.FindAllIndex([]byte(s), -1)

	if hyphens == nil {
		return
	}

	for _, loc := range hyphens {
		locStart := loc[0]
		t1 := s[0:locStart]

		sa.AddString(t1)
	}

	return
}
