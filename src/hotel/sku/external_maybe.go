package sku

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type ExternalMaybeLike interface {
	GetKennungLike() kennung.Kennung
	GetFDs() ExternalFDs
	kennung.FDPairGetter
}

type ExternalMaybe struct {
	Kennung kennung.KennungPtr
	FDs     ExternalFDs
}

func (a ExternalMaybe) String() string {
	return a.Kennung.String()
}

func (a ExternalMaybe) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a ExternalMaybe) Equals(b ExternalMaybe) bool {
	if a.Kennung.String() != b.Kennung.String() {
		return false
	}

	if !a.FDs.Equals(b.FDs) {
		return false
	}

	return true
}

func (e ExternalMaybe) GetId() schnittstellen.ValueLike {
	return e.Kennung
}

func (e ExternalMaybe) GetKennungLike() kennung.Kennung {
	return e.Kennung
}

func (e *ExternalMaybe) GetKennungLikePtr() kennung.KennungPtr {
	return e.Kennung
}

func (e ExternalMaybe) GetFDs() ExternalFDs {
	return e.FDs
}

func (e ExternalMaybe) GetObjekteFD() kennung.FD {
	return e.FDs.Objekte
}

func (e ExternalMaybe) GetAkteFD() kennung.FD {
	return e.FDs.Akte
}
