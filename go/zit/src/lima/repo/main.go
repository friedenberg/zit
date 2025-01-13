package repo

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

type Repo interface {
	GetRepoLayout() repo_layout.Layout
}

// TODO explore permissions for who can read / write from the archive
type Archive interface {
	Repo
	GetBlobStore() interfaces.BlobStore
	GetInventoryListStore() sku.InventoryListStore
}

type WorkingCopy interface {
	Archive

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

	// 	PullQueryGroupFromRemote2(
	// 		remote ReadWrite,
	// 		options RemoteTransferOptions,
	// 		query ...string,
	// 	) (err error)

	PullQueryGroupFromRemote(
		remote WorkingCopy,
		qg *query.Group,
		options RemoteTransferOptions,
	) (err error)

	ReadObjectHistory(
		oid *ids.ObjectId,
	) (skus []*sku.Transacted, err error)

	GetBlobStore() interfaces.BlobStore
}
