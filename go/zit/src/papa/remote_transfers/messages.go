package remote_transfers

import (
	"code.linenisgreat.com/zit-go/src/charlie/gattung"
	"code.linenisgreat.com/zit-go/src/charlie/sha"
	"code.linenisgreat.com/zit-go/src/india/matcher"
)

type messageRequestSkus struct {
	MetaSet matcher.Query
}

type messageRequestObjekteData struct {
	Gattung gattung.Gattung
	Sha     sha.Sha
}
