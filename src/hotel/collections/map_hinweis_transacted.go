package collections

import (
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
)

type MapHinweisTransacted map[hinweis.Hinweis]stored_zettel.Transacted

func MakeMapHinweisTransacted() MapHinweisTransacted {
	return make(MapHinweisTransacted)
}

func (m *MapHinweisTransacted) Add(z stored_zettel.Transacted) {
	(*m)[z.Hinweis] = z
}

func (m MapHinweisTransacted) Get(h hinweis.Hinweis) (z stored_zettel.Transacted, ok bool) {
	z, ok = m[h]
	return
}

func (a MapHinweisTransacted) Merge(b MapHinweisTransacted) {
	for _, z := range b {
		a.Add(z)
	}
}

func (a MapHinweisTransacted) Contains(z stored_zettel.Transacted) bool {
	_, ok := a[z.Hinweis]
	return ok
}

func (m MapHinweisTransacted) ToSlice() (s SliceTransacted) {
	s = make(SliceTransacted, 0, len(m))

	for _, sz := range m {
		s = append(s, sz)
	}

	return
}
