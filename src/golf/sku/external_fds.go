package sku

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
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

func (e ExternalFDs) GetCheckoutMode() (m checkout_mode.Mode, err error) {
	switch {
	case !e.Objekte.IsEmpty() && !e.Akte.IsEmpty():
		m = checkout_mode.ModeObjekteAndAkte

	case !e.Akte.IsEmpty():
		m = checkout_mode.ModeAkteOnly

	case !e.Objekte.IsEmpty():
		m = checkout_mode.ModeObjekteOnly

	default:
		err = checkout_mode.MakeErrInvalidCheckoutMode(
			errors.Errorf("all FD's are empty"),
		)
	}

	return
}
