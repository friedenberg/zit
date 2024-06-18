package query

import (
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type EtikettLua struct {
	*Lua
	*Kennung
}

func (k *EtikettLua) ContainsSku(sk *sku.Transacted) bool {
	if k.Kennung.ContainsSku(sk) {
		return true
	}

	if k.Lua.ContainsSku(sk) {
		return true
	}

	return false
}

func (k *EtikettLua) String() string {
	return k.Kennung.String()
}
