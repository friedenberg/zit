package etikett

type ExpanderAll struct{}

func (ex ExpanderAll) Expand(e Etikett) (out Set) {
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
