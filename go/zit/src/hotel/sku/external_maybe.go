package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

type ExternalMaybe struct {
	Kennung kennung.Kennung2
	FDs     ExternalFDs
}

func (a *ExternalMaybe) String() string {
	return a.Kennung.String()
}

func (a *ExternalMaybe) Equals(b ExternalMaybe) bool {
	if a.Kennung.String() != b.Kennung.String() {
		return false
	}

	if !a.FDs.Equals(&b.FDs) {
		return false
	}

	return true
}

func (e *ExternalMaybe) GetKennungLike() kennung.Kennung {
	return &e.Kennung
}

func (e *ExternalMaybe) GetKennungLikePtr() kennung.KennungPtr {
	return &e.Kennung
}

func (e *ExternalMaybe) GetFDs() *ExternalFDs {
	return &e.FDs
}

func (e *ExternalMaybe) GetObjekteFD() *fd.FD {
	return &e.FDs.Objekte
}

func (e *ExternalMaybe) GetAkteFD() *fd.FD {
	return &e.FDs.Akte
}
