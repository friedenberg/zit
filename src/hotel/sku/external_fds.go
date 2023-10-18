package sku

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/echo/fd"
)

type ExternalFDs struct {
	Objekte fd.FD
	Akte    fd.FD
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

func (e ExternalFDs) MakeConflictMarker() (path string) {
	path = fmt.Sprintf("%s.conflict", e.Objekte.GetPath())

	return
}

func (e ExternalFDs) conflictMarkerExists(fd *fd.FD) (ok bool) {
	if files.Exists(fmt.Sprintf("%s.conflict", fd)) {
		ok = true
	}

	return
}

func (e *ExternalFDs) ConflictMarkerError() (err error) {
	if e.conflictMarkerExists(&e.Objekte) {
		err = ErrExternalHasConflictMarker
	}

	if e.conflictMarkerExists(&e.Akte) {
		err = ErrExternalHasConflictMarker
	}

	return
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
