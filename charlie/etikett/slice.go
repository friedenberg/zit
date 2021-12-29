package etikett

type Slice struct {
	idx    int
	Values []Etikett
}

func (s Slice) Len() int {
	return len(s.Values)
}

func (s *Slice) SetIndex(i int) {
	s.idx = i
}

func (s Slice) CouldBe(r rune) bool {
	if r >= 'a' && r <= 'z' {
		return true
	}

	if r >= '0' && r <= '9' {
		return true
	}

	if r == '_' || r == '-' {
		return true
	}

	return false
}

func (s Slice) String() string {
	return s.Values[s.idx].String()
}

func (s Slice) Set(v string) (err error) {
	for s.Len() < s.idx+1 {
		s.Values = append(
			s.Values,
			Etikett{},
		)
	}

	if err = s.Values[s.idx].Set(v); err != nil {
		err = _Error(err)
		return
	}

	return
}
