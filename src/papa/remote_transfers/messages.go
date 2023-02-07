package remote_transfers

import (
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type messageRequestSkus struct {
	Filter       kennung.Filter
	GattungSlice []gattung.Gattung
	Sigil        kennung.Sigil
}

type messageRequestObjekteData struct {
	Gattung gattung.Gattung
	Sha     sha.Sha
}
