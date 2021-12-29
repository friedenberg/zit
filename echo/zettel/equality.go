package zettel

import "log"

func (z Zettel) Equals(z1 Zettel) bool {
	if !z.Akte.Equals(z1.Akte) {
		log.Print("akte")
		log.Print(z.Akte)
		log.Print(z1.Akte)
		return false
	}

	if !z.AkteExt.Equals(z1.AkteExt) {
		log.Print("akteext")
		return false
	}

	if z.Bezeichnung != z1.Bezeichnung {
		log.Print("bezeichnung")
		return false
	}

	if !z.Etiketten.Equals(z1.Etiketten) {
		log.Print("etiketten")
		return false
	}

	return true
}
