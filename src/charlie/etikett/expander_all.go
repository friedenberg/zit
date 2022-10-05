package etikett

type ExpanderAll struct{}

func (ex ExpanderAll) Expand(e Etikett) (expanded Set) {
	expanded = Set(newSetExpanded())
  expanded.open()
  defer expanded.close()

	expanded.addOnlyExact(e)

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

		expanded.addOnlyExact(Etikett{Value: t1})
		expanded.addOnlyExact(Etikett{Value: t2})

		if 0 < i && i < len(hyphens) {
			t1 := s[prevLocEnd:locStart]
			expanded.addOnlyExact(Etikett{Value: t1})
		}

		prevLocEnd = locEnd
	}

	return
}
