package kennung_index

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type EtikettVerzeichnisse struct {
	SortedExpanded []string
	Sorted         []string
}

func (z *EtikettVerzeichnisse) ResetWithEtikettSet(es kennung.EtikettSet) {
	ex := kennung.Expanded(es, kennung.ExpanderAll)
	z.SortedExpanded = ex.SortedString()
	z.Sorted = es.SortedString()
}

func (z *EtikettVerzeichnisse) Reset() {
	z.SortedExpanded = []string{}
	z.Sorted = []string{}

	z.SortedExpanded = z.SortedExpanded[:0]
	z.Sorted = z.Sorted[:0]
}

func (z *EtikettVerzeichnisse) ResetWith(z1 EtikettVerzeichnisse) {
	errors.TodoP4("improve performance by reusing slices")

	z.SortedExpanded = make([]string, len(z1.SortedExpanded))
	copy(z.SortedExpanded, z1.SortedExpanded)

	z.Sorted = make([]string, len(z1.Sorted))
	copy(z.Sorted, z1.Sorted)
}
