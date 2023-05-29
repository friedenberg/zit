package typen_index

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type Index interface {
	DidRead() bool
	HasChanges() bool
	Reset() error
	ExpandTyp(kennung.Typ) (Indexed, bool)
	StoreTyp(kennung.Typ) (err error)
	io.WriterTo
	io.ReaderFrom
}

type Indexed interface {
	GetTyp() kennung.Typ
	GetTridex() schnittstellen.Tridex
	GetTypenExpandedRight() schnittstellen.Set[kennung.Typ]
	GetTypenExpandedAll() schnittstellen.Set[kennung.Typ]
}

func MakeIndex() Index {
	return makeIndex()
}
