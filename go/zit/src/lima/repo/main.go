package repo

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

type Repo interface {
	GetRepo() Repo

	MakeQueryGroup(
		metaBuilder any,
		repoId ids.RepoId,
		externalQueryOptions sku.ExternalQueryOptions,
		args ...string,
	) (qg *query.Group, err error)

	MakeInventoryList(
		qg *query.Group,
	) (list *sku.List, err error)

	PullQueryGroupFromRemote(
		remote Repo,
		qg *query.Group,
		printCopies bool,
	) (err error)

	ReadObjectHistory(
		oid *ids.ObjectId,
	) (skus []*sku.Transacted, err error)

	GetBlobStore() dir_layout.BlobStore
}
