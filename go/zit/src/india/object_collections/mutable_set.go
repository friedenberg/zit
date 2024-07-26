package object_collections

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type MutableSet = interfaces.MutableSetLike[sku.ExternalLike]

func MakeMutableSetUniqueHinweis(zs ...sku.ExternalLike) MutableSet {
	return collections_value.MakeMutableValueSet(
		ids.ObjectIdKeyer[sku.ExternalLike]{},
		zs...,
	)
}

func MakeMutableSetUniqueFD(zs ...sku.ExternalLike) MutableSet {
	return collections_value.MakeMutableValueSet(
		KeyerFD{},
		zs...,
	)
}

func MakeMutableSetUniqueBlob(zs ...sku.ExternalLike) MutableSet {
	return collections_value.MakeMutableValueSet(
		KeyerBlob{},
		zs...,
	)
}
