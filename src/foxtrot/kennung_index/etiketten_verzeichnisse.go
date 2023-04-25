package kennung_index

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/tridex"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type EtikettenVerzeichnisse struct {
	Tridex         schnittstellen.MutableTridex
	Etiketten      schnittstellen.Set[kennung.Etikett]
	Expanded       schnittstellen.Set[kennung.Etikett]
	SortedExpanded []string
	Sorted         []string
}

func (z EtikettenVerzeichnisse) GetEtiketten() schnittstellen.Set[kennung.Etikett] {
	return z.Etiketten
}

func (z EtikettenVerzeichnisse) GetEtikettenExpanded() schnittstellen.Set[kennung.Etikett] {
	return z.Expanded
}

func (z *EtikettenVerzeichnisse) ResetWithEtikettSet(es kennung.EtikettSet) {
	if es == nil {
		es = kennung.MakeEtikettSet()
	}

	ex := kennung.Expanded(es, kennung.ExpanderAll)
	z.Tridex = tridex.Make(collections.SortedStrings[kennung.Etikett](ex)...)
	z.Etiketten = es.ImmutableClone()
	z.Expanded = kennung.Expanded(es, kennung.ExpanderRight)
	z.SortedExpanded = collections.SortedStrings[kennung.Etikett](ex)
	z.Sorted = collections.SortedStrings[kennung.Etikett](es)
}

func (z *EtikettenVerzeichnisse) Reset() {
	z.SortedExpanded = []string{}
	z.Sorted = []string{}

	z.SortedExpanded = z.SortedExpanded[:0]
	z.Etiketten = kennung.MakeEtikettSet()
	z.Expanded = kennung.MakeEtikettSet()
	z.Sorted = z.Sorted[:0]
	z.Tridex = tridex.Make()
}

func (z *EtikettenVerzeichnisse) ResetWith(z1 EtikettenVerzeichnisse) {
	errors.TodoP4("improve performance by reusing slices")

	z.Tridex = z1.Tridex.MutableClone()

	z.Etiketten = z1.Etiketten.ImmutableClone()
	z.Expanded = z1.Expanded.ImmutableClone()

	z.SortedExpanded = make([]string, len(z1.SortedExpanded))
	copy(z.SortedExpanded, z1.SortedExpanded)

	z.Sorted = make([]string, len(z1.Sorted))
	copy(z.Sorted, z1.Sorted)
}
