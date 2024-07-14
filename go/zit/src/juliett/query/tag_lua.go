package query

import (
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type TagLua struct {
	*Lua
	*ObjectId
}

func (k *TagLua) ContainsSku(sk *sku.Transacted) bool {
	if k.ObjectId.ContainsSku(sk) {
		return true
	}

	if k.Lua.ContainsSku(sk) {
		return true
	}

	return false
}

func (k *TagLua) String() string {
	return k.ObjectId.String()
}
