package etiketten_index

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/tridex"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type indexed struct {
	Etikett       kennung.Etikett
	ExpandedAll   schnittstellen.Set[kennung.Etikett]
	ExpandedRight schnittstellen.Set[kennung.Etikett]
	Tridex        schnittstellen.MutableTridex
}

func makeIndexed(k kennung.Etikett) (i *indexed) {
	i = &indexed{}
	i.ResetWithEtikett(k)
	return
}

func (i *indexed) ResetWithEtikett(k kennung.Etikett) {
	i.Etikett = k
	i.ExpandedAll = kennung.ExpandOne(k, kennung.ExpanderAll)
	i.ExpandedRight = kennung.ExpandOne(k, kennung.ExpanderRight)
	i.Tridex = tridex.Make(collections.SortedStrings[kennung.Etikett](i.ExpandedRight)...)
}

func (z indexed) GetEtikett() kennung.Etikett {
	return z.Etikett
}

func (z indexed) GetTridex() schnittstellen.Tridex {
	return z.Tridex
}

func (z indexed) GetEtikettenExpandedRight() schnittstellen.Set[kennung.Etikett] {
	return z.ExpandedRight
}

func (z indexed) GetEtikettenExpandedAll() schnittstellen.Set[kennung.Etikett] {
	return z.ExpandedAll
}

func (z *indexed) Reset() {
	z.Etikett.Reset()
	z.ExpandedRight = kennung.MakeEtikettSet()
	z.ExpandedAll = kennung.MakeEtikettSet()
	z.Tridex = tridex.Make()
}

func (z *indexed) ResetWith(z1 indexed) {
	z.ExpandedRight = z1.ExpandedRight.ImmutableClone()
	z.ExpandedAll = z1.ExpandedAll.ImmutableClone()
	z.Etikett = z1.Etikett
	z.Tridex = z1.Tridex.MutableClone()
}
