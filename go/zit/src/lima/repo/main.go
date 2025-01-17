package repo

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

// TODO explore permissions for who can read / write from the inventory list
// store
type Repo interface {
	GetEnv() env_ui.Env
	GetStoreVersion() interfaces.StoreVersion
	GetRepoType() repo_type.Type
	GetBlobStore() interfaces.BlobStore
	GetInventoryListStore() sku.InventoryListStore
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
