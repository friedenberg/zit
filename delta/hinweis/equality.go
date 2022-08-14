package hinweis

func (a Hinweis) Equals(b Hinweis) bool {
	if a.String() != b.String() {
		return false
	}

	return true
}
