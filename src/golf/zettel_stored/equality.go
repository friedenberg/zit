package zettel_stored

import "github.com/friedenberg/zit/src/alfa/errors"

func (a Stored) Equals(b Stored) bool {
	if !a.Zettel.Equals(b.Zettel) {
		errors.Print("zettel")
		return false
	}

	if !a.Sha.Equals(b.Sha) {
		errors.Print("sha")
		return false
	}

	return true
}
