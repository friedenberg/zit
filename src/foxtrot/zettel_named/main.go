package zettel_named

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/echo/zettel_stored"
)

type Zettel struct {
	Stored  zettel_stored.Stored
	Hinweis hinweis.Hinweis
}

func (a Zettel) Equals(b Zettel) bool {
	if !a.Stored.Equals(b.Stored) {
		errors.Print("stored")
		return false
	}

	if !a.Hinweis.Equals(b.Hinweis) {
		errors.Print("hinweis")
		return false
	}

	return true
}
