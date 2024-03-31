package query

import (
	"sync"

	"code.linenisgreat.com/zit/src/hotel/sku"
)

type VirtualStoreInitable struct {
	sku.VirtualStore
	sync.Once
}

func (ve *VirtualStoreInitable) Initialize() (err error) {
	ve.Do(func() { err = ve.VirtualStore.Initialize() })
	return
}

type Virtual struct {
	sku.VirtualStore
	Kennung
}

func (ve *Virtual) ContainsSku(sk *sku.Transacted) bool {
	if !ve.VirtualStore.ContainsSku(sk) {
		return false
	}

	if !ve.Kennung.ContainsSku(sk) {
		return false
	}

	return true
}
