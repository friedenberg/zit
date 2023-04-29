package sku

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type ExternalMaybeLike interface {
	IdLikeGetter
	kennung.FDPairGetter
	checkout_mode.Getter
}

type ExternalMaybe[
	T kennung.KennungLike[T],
	T1 kennung.KennungLikePtr[T],
] struct {
	Kennung T
	FDs     ExternalFDs
}

func (a ExternalMaybe[T, T1]) String() string {
	return a.Kennung.String()
}

func (a ExternalMaybe[T, T1]) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a ExternalMaybe[T, T1]) Equals(b ExternalMaybe[T, T1]) bool {
	if !a.Kennung.Equals(b.Kennung) {
		return false
	}

	if !a.FDs.Equals(b.FDs) {
		return false
	}

	return true
}

func (e ExternalMaybe[T, T1]) GetId() schnittstellen.IdLike {
	return e.Kennung
}

func (e ExternalMaybe[T, T1]) GetIdLike() schnittstellen.IdLike {
	return e.Kennung
}

func (e ExternalMaybe[T, T1]) GetObjekteFD() kennung.FD {
	return e.FDs.Objekte
}

func (e ExternalMaybe[T, T1]) GetAkteFD() kennung.FD {
	return e.FDs.Akte
}

func (e ExternalMaybe[T, T1]) GetCheckoutMode() (m checkout_mode.Mode, err error) {
	switch {
	case !e.FDs.Objekte.IsEmpty() && !e.FDs.Akte.IsEmpty():
		m = checkout_mode.ModeObjekteAndAkte

	case !e.FDs.Akte.IsEmpty():
		m = checkout_mode.ModeAkteOnly

	case !e.FDs.Objekte.IsEmpty():
		m = checkout_mode.ModeObjekteOnly

	default:
		err = MakeErrInvalidCheckoutMode(
			errors.Errorf("all FD's are empty"),
		)
	}

	return
}
