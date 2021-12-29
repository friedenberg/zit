package stored_zettel

import "log"

func (a Stored) Equals(b Stored) bool {
	if !a.Zettel.Equals(b.Zettel) {
		log.Print("zettel")
		return false
	}

	if !a.Mutter.Equals(b.Mutter) {
		log.Print("mutter")
		return false
	}

	if !a.Kinder.Equals(b.Kinder) {
		log.Print("sha")
		return false
	}

	if !a.Sha.Equals(b.Sha) {
		log.Print("sha")
		return false
	}

	return true
}

func (a Named) Equals(b Named) bool {
	if !a.Stored.Equals(b.Stored) {
		log.Print("stored")
		return false
	}

	if !a.Hinweis.Equals(b.Hinweis) {
		log.Print("hinweis")
		return false
	}

	return true
}
