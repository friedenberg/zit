package stored_zettel

import (
	"github.com/friedenberg/zit/src/alfa/logz"
)

func (a Stored) Equals(b Stored) bool {
	if !a.Zettel.Equals(b.Zettel) {
		logz.Print("zettel")
		return false
	}

	if !a.Mutter.Equals(b.Mutter) {
		logz.Print("mutter")
		return false
	}

	if !a.Kinder.Equals(b.Kinder) {
		logz.Print("sha")
		return false
	}

	if !a.Sha.Equals(b.Sha) {
		logz.Print("sha")
		return false
	}

	return true
}

func (a Named) Equals(b Named) bool {
	if !a.Stored.Equals(b.Stored) {
		logz.Print("stored")
		return false
	}

	if !a.Hinweis.Equals(b.Hinweis) {
		logz.Print("hinweis")
		return false
	}

	return true
}
