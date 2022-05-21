package organize_text

import (
	"io"

	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

type Text interface {
	io.ReaderFrom
	io.WriterTo
	Etiketten() etikett.Set
	ZettelsExisting() map[string]zettelSet
	ZettelsNew() map[string]newZettelSet
	ChangesFrom(Text) Changes
}

type organizeText struct {
	etiketten etikett.Set
	zettels   assignments
}

func (t organizeText) Etiketten() etikett.Set {
	return t.etiketten
}

func (t organizeText) ZettelsExisting() map[string]zettelSet {
	return t.zettels.etikettenToExisting
}

func (t organizeText) ZettelsNew() map[string]newZettelSet {
	return t.zettels.etikettenToNew
}

func New(options Options, named stored_zettel.SetNamed) (ot *organizeText, err error) {
	ot = NewEmpty()

	ot.etiketten = options.RootEtiketten

	for _, z := range named {
		groups := options.GroupZettel(z)

		for _, g := range groups {
			ot.zettels.AddStored(g.String(), z)
		}
	}

	return
}

func NewEmpty() (ot *organizeText) {
	ot = &organizeText{
		zettels: newEtikettToZettels(),
	}

	return
}
