package typen_index

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/kennung_index"
)

type Index = kennung_index.Index2[kennung.Typ]

func MakeIndex() Index {
	return kennung_index.MakeIndex2[kennung.Typ]()
}
