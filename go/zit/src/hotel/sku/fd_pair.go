package sku

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/echo/thyme"
)

var ErrExternalHasConflictMarker = errors.New("external has conflict marker")

type FDPair struct {
	Objekte fd.FD
	Akte    fd.FD
}

func (ef *FDPair) GetTai() kennung.Tai {
	return kennung.TaiFromTime(ef.LatestModTime())
}

func (ef *FDPair) GetTime() thyme.Time {
	return ef.LatestModTime()
}

func (ef *FDPair) LatestModTime() thyme.Time {
	o, a := ef.Objekte.ModTime(), ef.Akte.ModTime()

	if o.Less(a) {
		return a
	} else {
		return o
	}
}

// func (dst *ExternalFDs) String() string {
//   panic("")
// }

func (dst *FDPair) ResetWith(src *FDPair) {
	dst.Objekte.ResetWith(&src.Objekte)
	dst.Akte.ResetWith(&src.Akte)
}

func (a *FDPair) Equals(b *FDPair) bool {
	if !a.Objekte.Equals(&b.Objekte) {
		return false
	}

	if !a.Akte.Equals(&b.Akte) {
		return false
	}

	return true
}

func (e *FDPair) MakeConflictMarker() (path string) {
	path = fmt.Sprintf("%s.conflict", e.Objekte.GetPath())

	return
}

func (e *FDPair) conflictMarkerExists(fd *fd.FD) (ok bool) {
	if files.Exists(fmt.Sprintf("%s.conflict", fd)) {
		ok = true
	}

	return
}

func (e *FDPair) ConflictMarkerError() (err error) {
	if e.conflictMarkerExists(&e.Objekte) {
		err = errors.Wrapf(ErrExternalHasConflictMarker, "Objekte: %s", &e.Objekte)
		return
	}

	if e.conflictMarkerExists(&e.Akte) {
		err = errors.Wrapf(ErrExternalHasConflictMarker, "Akte: %s", &e.Akte)
		return
	}

	return
}

func (e *FDPair) GetCheckoutModeOrError() (m checkout_mode.Mode, err error) {
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

func (e *FDPair) GetCheckoutMode() (m checkout_mode.Mode) {
	switch {
	case !e.Objekte.IsEmpty() && !e.Akte.IsEmpty():
		m = checkout_mode.ModeObjekteAndAkte

	case !e.Akte.IsEmpty():
		m = checkout_mode.ModeAkteOnly

	case !e.Objekte.IsEmpty():
		m = checkout_mode.ModeObjekteOnly
	}

	return
}
