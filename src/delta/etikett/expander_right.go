package etikett

type ExpanderRight struct{}

func (ex ExpanderRight) Expand(e Etikett) (expanded Set) {
	expanded = Set(newSetExpanded())
	expanded.addOnlyExact(e)

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

		expanded.addOnlyExact(Etikett{Value: t1})
	}

	return
}
