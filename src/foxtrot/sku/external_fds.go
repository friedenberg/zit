package sku

import (
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type ExternalFDs struct {
	Objekte kennung.FD
	Akte    kennung.FD
}

func (a ExternalFDs) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a ExternalFDs) Equals(b ExternalFDs) bool {
	if !a.Objekte.Equals(b.Objekte) {
		return false
	}

	if !a.Akte.Equals(b.Akte) {
		return false
	}

	return true
}
