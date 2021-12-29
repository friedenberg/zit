package etikett

import "regexp"

var (
	regexExpandTagsHyphens *regexp.Regexp
)

func init() {
	regexExpandTagsHyphens = regexp.MustCompile(`-`)
}

type Expander interface {
	Expand(Etikett) Set
}

func (e Etikett) Expanded(exes ...Expander) (expanded Set) {
	expanded = NewSet()

	if len(exes) == 0 {
		exes = []Expander{ExpanderAll{}}
	}

	for _, ex := range exes {
		for _, e := range ex.Expand(e) {
			expanded.Add(e)
		}
	}

	return
}

type ExpanderAll struct{}

func (ex ExpanderAll) Expand(e Etikett) (expanded Set) {
	expanded = NewSet()
	expanded.Add(e)

	s := e.String()

	if s == "" {
		return
	}

	hyphens := regexExpandTagsHyphens.FindAllIndex([]byte(s), -1)

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

		expanded.Add(Etikett{Value: t1})
		expanded.Add(Etikett{Value: t2})

		if 0 < i && i < len(hyphens) {
			t1 := s[prevLocEnd:locStart]
			expanded.Add(Etikett{Value: t1})
		}

		prevLocEnd = locEnd
	}

	return
}

type ExpanderRight struct{}

func (ex ExpanderRight) Expand(e Etikett) (expanded Set) {
	expanded = NewSet()
	expanded.Add(e)

	s := e.String()

	if s == "" {
		return
	}

	hyphens := regexExpandTagsHyphens.FindAllIndex([]byte(s), -1)

	if hyphens == nil {
		return
	}

	for _, loc := range hyphens {
		locStart := loc[0]
		t1 := s[0:locStart]

		expanded.Add(Etikett{Value: t1})
	}

	return
}
