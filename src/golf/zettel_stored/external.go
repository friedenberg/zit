package zettel_stored

import "fmt"

func (e External) String() string {
	return e.ExternalPathAndSha()
}

func (e External) ExternalPathAndSha() string {
	if !e.ZettelFD.IsEmpty() {
		return fmt.Sprintf("[%s %s]", e.ZettelFD.Path, e.Stored.Sha)
	} else if !e.AkteFD.IsEmpty() {
		return fmt.Sprintf("[%s %s]", e.AkteFD.Path, e.Stored.Zettel.Akte)
	} else {
		return ""
	}
}
