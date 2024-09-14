package query

import (
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type TagLua struct {
	*Lua
	*ObjectId
}

func (k *TagLua) ContainsSku(tg sku.TransactedGetter) bool {
	if k.ObjectId.ContainsSku(tg) {
		return true
	}

	if k.Lua.ContainsSku(tg) {
		return true
	}

	return false
}

func (k *TagLua) String() string {
	return k.ObjectId.String()
}
