package sku

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type ExternalMaybeLike interface {
	GetKennungLike() kennung.Kennung
	kennung.FDPairGetter
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

func (e ExternalMaybe[T, T1]) GetId() schnittstellen.ValueLike {
	return e.Kennung
}

func (e ExternalMaybe[T, T1]) GetKennungLike() kennung.Kennung {
	return e.Kennung
}

func (e *ExternalMaybe[T, T1]) GetKennungPtr() kennung.KennungPtr {
	return T1(&e.Kennung)
}

func (e ExternalMaybe[T, T1]) GetFDs() ExternalFDs {
	return e.FDs
}

func (e ExternalMaybe[T, T1]) GetObjekteFD() kennung.FD {
	return e.FDs.Objekte
}

func (e ExternalMaybe[T, T1]) GetAkteFD() kennung.FD {
	return e.FDs.Akte
}
