package organize_text

import (
	"io"
)

type Text interface {
	io.ReaderFrom
	io.WriterTo
	Etiketten() _EtikettSet
	ZettelsExisting() map[string]zettelSet
	ZettelsNew() map[string]newZettelSet
	ChangesFrom(Text) Changes
}

type organizeText struct {
	etiketten _EtikettSet
	zettels   assignments
}

func (t organizeText) Etiketten() _EtikettSet {
	return t.etiketten
}

func (t organizeText) ZettelsExisting() map[string]zettelSet {
	return t.zettels.etikettenToExisting
}

func (t organizeText) ZettelsNew() map[string]newZettelSet {
	return t.zettels.etikettenToNew
}

func New(options Options, zettels map[string]_NamedZettel) (ot *organizeText, err error) {
	ot = NewEmpty()

	ot.etiketten = options.RootEtiketten

	for _, z := range zettels {
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
