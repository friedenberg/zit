package zettel

import (
	"github.com/friedenberg/zit/src/alfa/logz"
)

func (z Zettel) Equals(z1 Zettel) bool {
	if !z.Akte.Equals(z1.Akte) {
		logz.Print("akte")
		logz.Print(z.Akte)
		logz.Print(z1.Akte)
		return false
	}

	if !z.AkteExt.Equals(z1.AkteExt) {
		logz.Print("akteext")
		return false
	}

	if z.Bezeichnung != z1.Bezeichnung {
		logz.Print("bezeichnung")
		return false
	}

	if !z.Etiketten.Equals(z1.Etiketten) {
		logz.Print("etiketten")
		return false
	}

	return true
}
