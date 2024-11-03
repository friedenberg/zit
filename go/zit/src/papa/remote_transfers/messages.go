package remote_transfers

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

type messageRequestSkus struct {
	MetaSet *query.Group
}

type messageRequestObjectData struct {
	Gattung genres.Genre
	Sha     sha.Sha
}
