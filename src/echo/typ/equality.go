package typ

import "strings"

func (a Kennung) Equals(b Kennung) bool {
	if a.String() != b.String() {
		return false
	}

	return true
}

func (a Kennung) Contains(b Kennung) bool {
	as := a.String()

	if as == "" {
		return true
	}

	if !strings.HasPrefix(b.String(), as) {
		return false
	}

	return true
}
