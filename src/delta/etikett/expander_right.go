package etikett

type ExpanderRight struct{}

func (ex ExpanderRight) Expand(e Etikett) (out Set) {
	expanded := MakeMutableSet()

	defer func() {
		out = Set(newSetExpanded(expanded.Elements()...))
	}()

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
