package etiketten_index

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type Index interface {
	ExpandEtikett(kennung.Etikett) (Indexed, bool)
	StoreEtikett(kennung.Etikett) (err error)
	io.WriterTo
	io.ReaderFrom
}

type Indexed interface {
	GetEtikett() kennung.Etikett
	GetTridex() schnittstellen.Tridex
	GetEtikettenExpandedRight() schnittstellen.Set[kennung.Etikett]
	GetEtikettenExpandedAll() schnittstellen.Set[kennung.Etikett]
}

func MakeIndex() Index {
	return makeIndex()
}
