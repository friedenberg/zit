package external_store

import "code.linenisgreat.com/zit/go/zit/src/hotel/sku"

type (
	SkuType           = sku.ExternalLike
	SkuTypeSet        = sku.ExternalLikeSet
	SkuTypeSetMutable = sku.ExternalLikeMutableSet
	ObjectFactory     = objectFactoryExternalLike
)

var (
	MakeSkuType                = makeExternalLike
	CloneSkuType               = cloneExternalLike
	CloneSkuTypeFromTransacted = cloneFromTransactedExternalLike
	MakeSkuTypeSetMutable      = sku.MakeExternalLikeMutableSet
)

// type (
// 	SkuType           = *sku.CheckedOut
// 	SkuTypeSet        = sku.CheckedOutSet
// 	SkuTypeSetMutable = sku.CheckedOutMutableSet
// 	ObjectFactory     = objectFactoryCheckedOut
// )

// var (
// 	MakeSkuType                = makeCheckedOut
// 	CloneSkuType               = cloneCheckedOut
// 	CloneSkuTypeFromTransacted = cloneFromTransactedCheckedOut
// 	MakeSkuTypeSetMutable      = sku.MakeCheckedOutMutableSet
// )
