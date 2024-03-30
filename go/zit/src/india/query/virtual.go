package query

import (
	"sync"

	"code.linenisgreat.com/zit/src/hotel/sku"
)

type VirtualStore interface {
	Init() error
	ContainsMatchable(*sku.Transacted) bool
}

type VirtualStoreInitable struct {
	VirtualStore
	sync.Once
}

func (ve *VirtualStoreInitable) Init() (err error) {
	ve.Do(func() { err = ve.VirtualStore.Init() })
	return
}

type Virtual struct {
	VirtualStore
	Kennung
}

func (ve *Virtual) ContainsMatchable(sk *sku.Transacted) bool {
	if !ve.VirtualStore.ContainsMatchable(sk) {
		return false
	}

	if !ve.Kennung.ContainsMatchable(sk) {
		return false
	}

	return true
}
