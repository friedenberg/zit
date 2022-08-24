package typ

func (a Typ) Equals(b Typ) bool {
	if a.String() != b.String() {
		return false
	}

	return true
}
