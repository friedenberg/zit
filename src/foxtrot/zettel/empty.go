package zettel

import "strings"

func (z Zettel) IsEmpty() bool {
	if strings.TrimSpace(z.Bezeichnung.String()) != "" {
		return false
	}

	if z.Etiketten.Len() > 0 {
		return false
	}

	if !z.Akte.IsNull() {
		return false
	}

	return true
}
