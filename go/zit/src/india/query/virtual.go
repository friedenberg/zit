package query

import (
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type VirtualStoreInitable struct {
	VirtualStore
	didInit  bool
	onceInit sync.Once
}

func (ve *VirtualStoreInitable) Initialize() (err error) {
	ve.onceInit.Do(func() {
		err = ve.VirtualStore.Initialize()
		ve.didInit = true
	})

	return
}

func (ve *VirtualStoreInitable) Flush() (err error) {
	if !ve.didInit {
		return
	}

	if err = ve.VirtualStore.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type Virtual struct {
	sku.Queryable
	Kennung
}

func (ve *VirtualStoreInitable) Query(
  qg *Group,
  f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	if err = ve.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = ve.VirtualStore.Query(qg, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (ve *Virtual) ContainsSku(sk *sku.Transacted) bool {
	if !ve.Queryable.ContainsSku(sk) {
		return false
	}

	if !ve.Kennung.ContainsSku(sk) {
		return false
	}

	return true
}
