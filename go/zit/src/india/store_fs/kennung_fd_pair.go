package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

type KennungFDPair struct {
	Kennung kennung.Kennung2
	FDs     FDPair
}

func (a *KennungFDPair) String() string {
	return a.Kennung.String()
}

func (a *KennungFDPair) Equals(b KennungFDPair) bool {
	if a.Kennung.String() != b.Kennung.String() {
		return false
	}

	if !a.FDs.Equals(&b.FDs) {
		return false
	}

	return true
}

func (e *KennungFDPair) GetKennungLike() kennung.Kennung {
	return &e.Kennung
}

func (e *KennungFDPair) GetKennungLikePtr() kennung.KennungPtr {
	return &e.Kennung
}

func (e *KennungFDPair) GetFDs() *FDPair {
	return &e.FDs
}

func (e *KennungFDPair) GetObjekteFD() *fd.FD {
	return &e.FDs.Objekte
}

func (e *KennungFDPair) GetAkteFD() *fd.FD {
	return &e.FDs.Akte
}
