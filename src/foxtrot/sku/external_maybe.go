package sku

import (
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type ExternalMaybeLike interface {
	GetId() IdLike
	kennung.FDPairGetter
}

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

type ExternalMaybe[T kennung.KennungLike[T], T1 kennung.KennungLikePtr[T]] struct {
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

func (e ExternalMaybe[T, T1]) GetId() IdLike {
	return e.Kennung
}

func (e ExternalMaybe[T, T1]) GetObjekteFD() kennung.FD {
	return e.FDs.Objekte
}

func (e ExternalMaybe[T, T1]) GetAkteFD() kennung.FD {
	return e.FDs.Akte
}
