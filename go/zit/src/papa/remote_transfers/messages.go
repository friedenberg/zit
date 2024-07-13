package remote_transfers

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

type messageRequestSkus struct {
	MetaSet *query.Group
}

type messageRequestObjekteData struct {
	Gattung gattung.Genre
	Sha     sha.Sha
}
