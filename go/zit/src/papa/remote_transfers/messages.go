package remote_transfers

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/india/matcher"
)

type messageRequestSkus struct {
	MetaSet matcher.Query
}

type messageRequestObjekteData struct {
	Gattung gattung.Gattung
	Sha     sha.Sha
}