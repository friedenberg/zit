package repo

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_immutable_io"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
)

// TODO explore permissions for who can read / write from the inventory list
// store
type Repo interface {
	GetEnv() env_ui.Env
	GetImmutableConfig() config_immutable_io.ConfigLoaded
	GetBlobStore() interfaces.BlobStore
	GetObjectStore() sku.ObjectStore
	GetTypedInventoryListBlobStore() typed_blob_store.InventoryList
	GetInventoryListStore() sku.InventoryListStore

	MakeImporter(
		options sku.ImporterOptions,
		storeOptions sku.StoreOptions,
	) sku.Importer

	// TODO switch to seq
	ImportList(
		list *sku.List,
		i sku.Importer,
	) (err error)
}

type WorkingCopy interface {
	Repo

	// MakeQueryGroup(
	// 	builderOptions query.BuilderOptions,
	// 	args ...string,
	// ) (qg *query.Group, err error)

	MakeExternalQueryGroup(
		builderOptions query.BuilderOptions,
		externalQueryOptions sku.ExternalQueryOptions,
		args ...string,
	) (qg *query.Group, err error)

	MakeInventoryList(
		qg *query.Group,
	) (list *sku.List, err error)

	PullQueryGroupFromRemote(
		remote Repo,
		qg *query.Group,
		options RemoteTransferOptions,
	) (err error)

	ReadObjectHistory(
		oid *ids.ObjectId,
	) (skus []*sku.Transacted, err error)
}
