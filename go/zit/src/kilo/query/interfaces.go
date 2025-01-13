package query

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

type (
	ExternalStoreForQuery interface {
		ParseAndApplyExternalObjectIdsForQuery(qg *Group, v string) error
		sku.ExternalStoreForQuery
	}

	ExternalStoreForQueryGetter interface {
		GetExternalStoreForQuery(ids.RepoId) (ExternalStoreForQuery, bool)
	}
)
