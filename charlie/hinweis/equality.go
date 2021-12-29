package hinweis

func (a hinweis) Equals(b Hinweis) bool {
	if a.String() != b.String() {
		return false
	}

	return true
}
