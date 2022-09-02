package zettel

import "github.com/friedenberg/zit/src/bravo/errors"

func (z Zettel) Equals(z1 Zettel) bool {
	if !z.Akte.Equals(z1.Akte) {
		errors.Print("akte")
		errors.Print(z.Akte)
		errors.Print(z1.Akte)
		return false
	}

	if !z.Typ.Equals(z1.Typ) {
		errors.Print("akteext")
		return false
	}

	if z.Bezeichnung != z1.Bezeichnung {
		errors.Print("bezeichnung")
		return false
	}

	if !z.Etiketten.Equals(z1.Etiketten) {
		errors.Print("etiketten")
		return false
	}

	return true
}
