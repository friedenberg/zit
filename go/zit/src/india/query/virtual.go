package query

import (
	"sync"

	"code.linenisgreat.com/zit/src/hotel/sku"
)

type VirtualStore interface {
	Initialize() error
	ContainsMatchable(*sku.Transacted) bool
}

type VirtualStoreInitable struct {
	VirtualStore
	sync.Once
}

func (ve *VirtualStoreInitable) Initialize() (err error) {
	ve.Do(func() { err = ve.VirtualStore.Initialize() })
	return
}

type Virtual struct {
	VirtualStore
	Kennung
}

func (ve *Virtual) ContainsSku(sk *sku.Transacted) bool {
	if !ve.VirtualStore.ContainsMatchable(sk) {
		return false
	}

	if !ve.Kennung.ContainsSku(sk) {
		return false
	}

	return true
}
