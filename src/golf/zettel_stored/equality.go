package zettel_stored

import (
	"github.com/friedenberg/zit/src/alfa/logz"
)

func (a Stored) Equals(b Stored) bool {
	if !a.Zettel.Equals(b.Zettel) {
		logz.Print("zettel")
		return false
	}

	if !a.Sha.Equals(b.Sha) {
		logz.Print("sha")
		return false
	}

	return true
}
