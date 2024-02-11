package remote_transfers

import (
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/india/matcher"
)

type messageRequestSkus struct {
	MetaSet matcher.Query
}

type messageRequestObjekteData struct {
	Gattung gattung.Gattung
	Sha     sha.Sha
}
