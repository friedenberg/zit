package remote_transfers

import (
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/hotel/matcher_proto"
)

type messageRequestSkus struct {
	MetaSet matcher_proto.QueryGroup
}

type messageRequestObjekteData struct {
	Gattung gattung.Gattung
	Sha     sha.Sha
}
