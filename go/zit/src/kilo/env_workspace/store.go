package env_workspace

// import "code.linenisgreat.com/zit/go/zit/src/juliett/sku"

// type ExternalStore interface {
// 	sku.ExternalStoreReadAllExternalItems
// 	sku.ExternalStoreUpdateTransacted
// 	sku.ExternalStoreReadExternalLikeFromObjectIdLike
// 	QueryCheckedOut
// }

type (
	StoreReadAllExternalItems interface {
		ReadAllExternalItems() error
	}
)
