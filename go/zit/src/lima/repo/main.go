package repo

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

type Repo interface {
	GetEnv() *env.Env
}

// TODO rename to Archive?
// TODO explore permissions for who can read / write from the archive
type Relay interface {
	Repo
	GetBlobStore() interfaces.BlobStore
}

// TODO rename to WorkingCopy?
type ReadWrite interface {
	Relay

	MakeQueryGroup(
		builderOptions query.BuilderOptions,
		repoId ids.RepoId,
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
		remote ReadWrite,
		qg *query.Group,
		options RemoteTransferOptions,
	) (err error)

	ReadObjectHistory(
		oid *ids.ObjectId,
	) (skus []*sku.Transacted, err error)

	GetBlobStore() interfaces.BlobStore
}
