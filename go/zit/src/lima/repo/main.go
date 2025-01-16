package repo

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

type Repo interface {
	GetRepoType() repo_type.Type
	// GetRepoLayout() repo_layout.Layout
}

// TODO explore permissions for who can read / write from the archive
type Archive interface {
	Repo
	GetBlobStore() interfaces.BlobStore
	GetInventoryListStore() sku.InventoryListStore
	// InitializeUnixSocket(config net.ListenConfig, path string) (UnixSocket, error)
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

	PullQueryGroupFromRemote(
		remote Archive,
		qg *query.Group,
		options RemoteTransferOptions,
	) (err error)

	ReadObjectHistory(
		oid *ids.ObjectId,
	) (skus []*sku.Transacted, err error)
}
