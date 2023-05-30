package etiketten_index

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/kennung_index"
)

type Index = kennung_index.Index2[kennung.Etikett]

func MakeIndex() Index {
	return kennung_index.MakeIndex2[kennung.Etikett]()
}
