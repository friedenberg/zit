package typ

import "strings"

func (a Typ) Equals(b Typ) bool {
	if a.String() != b.String() {
		return false
	}

	return true
}

func (a Typ) Contains(b Typ) bool {
	as := a.String()

	if as == "" {
		return true
	}

	if !strings.HasPrefix(b.String(), as) {
		return false
	}

	return true
}
