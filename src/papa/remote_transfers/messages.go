package remote_transfers

import (
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/foxtrot/id_set"
)

type messageRequestSkus struct {
	Filter       id_set.Filter
	GattungSlice []gattung.Gattung
}

type messageRequestObjekteData struct {
	Gattung gattung.Gattung
	Sha     sha.Sha
}
