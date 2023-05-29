package typen_index

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/tridex"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type indexed struct {
	Typ           kennung.Typ
	ExpandedAll   schnittstellen.Set[kennung.Typ]
	ExpandedRight schnittstellen.Set[kennung.Typ]
	Tridex        schnittstellen.MutableTridex
}

func makeIndexed(k kennung.Typ) (i *indexed) {
	i = &indexed{}
	i.ResetWithTyp(k)
	return
}

func (i *indexed) ResetWithTyp(k kennung.Typ) {
	i.Typ = k
	i.ExpandedAll = kennung.ExpandOne(k, kennung.ExpanderAll)
	i.ExpandedRight = kennung.ExpandOne(k, kennung.ExpanderRight)
	i.Tridex = tridex.Make(collections.SortedStrings[kennung.Typ](i.ExpandedRight)...)
}

func (z indexed) GetTyp() kennung.Typ {
	return z.Typ
}

func (z indexed) GetTridex() schnittstellen.Tridex {
	return z.Tridex
}

func (z indexed) GetTypenExpandedRight() schnittstellen.Set[kennung.Typ] {
	return z.ExpandedRight
}

func (z indexed) GetTypenExpandedAll() schnittstellen.Set[kennung.Typ] {
	return z.ExpandedAll
}

func (z *indexed) Reset() {
	z.Typ.Reset()
	z.ExpandedRight = kennung.MakeTypSet()
	z.ExpandedAll = kennung.MakeTypSet()
	z.Tridex = tridex.Make()
}

func (z *indexed) ResetWith(z1 indexed) {
	z.ExpandedRight = z1.ExpandedRight.ImmutableClone()
	z.ExpandedAll = z1.ExpandedAll.ImmutableClone()
	z.Typ = z1.Typ
	z.Tridex = z1.Tridex.MutableClone()
}
