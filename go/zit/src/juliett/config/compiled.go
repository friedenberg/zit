package config

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type compiled struct {
	lock sync.Locker

	changes []string

	Sku sku.Transacted

	mutable_config_private

	DefaultTags  ids.TagSet
	Tags         interfaces.MutableSetLike[*tag]
	ImplicitTags implicitTagMap

	// Typen
	ExtensionsToTypes map[string]string
	TypesToExtensions map[string]string
	DefaultType       sku.Transacted // deprecated
	Types             sku.TransactedMutableSet
	InlineTypes       interfaces.SetLike[values.String]

	// Kasten
	Repos sku.TransactedMutableSet
}

func (kc *compiled) IsInlineType(k ids.Type) (isInline bool) {
	todo.Change("fix this horrible hack")
	if k.IsEmpty() {
		return true
	}

	isInline = kc.InlineTypes.ContainsKey(k.String()) ||
		builtin_types.IsBuiltin(k)

	return
}
