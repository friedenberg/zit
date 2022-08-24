package collections

import (
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
)

type MapShaTransacted map[sha.Sha]stored_zettel.Transacted

func MakeMapShaTransacted() MapShaTransacted {
	return make(MapShaTransacted)
}

func (m *MapShaTransacted) Add(z stored_zettel.Transacted) {
	(*m)[z.Stored.Sha] = z
}

func (m MapShaTransacted) Get(s sha.Sha) (z stored_zettel.Transacted, ok bool) {
	z, ok = m[s]
	return
}

func (a MapShaTransacted) Merge(b MapShaTransacted) {
	for _, z := range b {
		a.Add(z)
	}
}

func (a MapShaTransacted) Contains(z stored_zettel.Transacted) bool {
	_, ok := a[z.Stored.Sha]
	return ok
}

func (m MapShaTransacted) ToSlice() (s SliceTransacted) {
	s = make(SliceTransacted, 0, len(m))

	for _, sz := range m {
		s = append(s, sz)
	}

	return
}
