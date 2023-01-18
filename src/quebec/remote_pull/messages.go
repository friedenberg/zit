package remote_pull

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/foxtrot/sha"
	"github.com/friedenberg/zit/src/golf/id_set"
)

type messageRequestSkus struct {
	Filter       id_set.Filter
	GattungSlice []gattung.Gattung
}

type messageRequestObjekteData struct {
	Gattung gattung.Gattung
	Sha     sha.Sha
}
