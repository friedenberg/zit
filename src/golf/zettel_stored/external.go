package zettel_stored

import "fmt"

func (e External) String() string {
	return e.ExternalPathAndSha()
}

func (e External) ExternalPathAndSha() string {
	if e.Path != "" {
		return fmt.Sprintf("[%s %s]", e.Path, e.Stored.Sha)
	} else if e.AktePath != "" {
		return fmt.Sprintf("[%s %s]", e.AktePath, e.Stored.Zettel.Akte)
	} else {
		return ""
	}
}
