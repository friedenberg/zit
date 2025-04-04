package store_workspace

import "code.linenisgreat.com/zit/go/zit/src/echo/ids"

type (
	Store interface {
		GetObjectIdsForString(string) ([]ids.ExternalObjectIdLike, error)
	}

	StoreGetter interface {
		GetWorkspaceStoreForQuery(ids.RepoId) (Store, bool)
	}
)
