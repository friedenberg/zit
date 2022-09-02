package zettel_named

import (
	"fmt"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/golf/zettel_stored"
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

func (z Zettel) String() string {
	return fmt.Sprintf("[%s %s]", z.Hinweis, z.Stored.Sha)
}
