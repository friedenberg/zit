package kennung_index

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/tridex"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type EtikettenVerzeichnisse struct {
	Tridex         *tridex.Tridex
	SortedExpanded []string
	Sorted         []string
}

func (z *EtikettenVerzeichnisse) ResetWithEtikettSet(es kennung.EtikettSet) {
	ex := kennung.Expanded(es, kennung.ExpanderAll)
	z.Tridex = tridex.Make(ex.SortedString()...)
	z.SortedExpanded = ex.SortedString()
	z.Sorted = es.SortedString()
}

func (z *EtikettenVerzeichnisse) Reset() {
	z.SortedExpanded = []string{}
	z.Sorted = []string{}

	z.SortedExpanded = z.SortedExpanded[:0]
	z.Sorted = z.Sorted[:0]
	z.Tridex = tridex.Make()
}

func (z *EtikettenVerzeichnisse) ResetWith(z1 EtikettenVerzeichnisse) {
	errors.TodoP4("improve performance by reusing slices")

	z.SortedExpanded = make([]string, len(z1.SortedExpanded))
	copy(z.SortedExpanded, z1.SortedExpanded)

	z.Sorted = make([]string, len(z1.Sorted))
	copy(z.Sorted, z1.Sorted)

	z.Tridex = z1.Tridex.Copy()
}
