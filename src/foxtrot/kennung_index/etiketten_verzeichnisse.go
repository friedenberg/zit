package kennung_index

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type EtikettenVerzeichnisse struct {
	Etiketten schnittstellen.Set[kennung.Etikett]
	Sorted    []string
}

func (z EtikettenVerzeichnisse) GetEtiketten() schnittstellen.Set[kennung.Etikett] {
	return z.Etiketten
}

func (z *EtikettenVerzeichnisse) ResetWithEtikettSet(es kennung.EtikettSet) {
	if es == nil {
		es = kennung.MakeEtikettSet()
	}

	z.Etiketten = es.ImmutableClone()
	z.Sorted = collections.SortedStrings[kennung.Etikett](es)
}

func (z *EtikettenVerzeichnisse) Reset() {
	z.Sorted = []string{}

	z.Etiketten = kennung.MakeEtikettSet()
	z.Sorted = z.Sorted[:0]
}

func (z *EtikettenVerzeichnisse) ResetWith(z1 EtikettenVerzeichnisse) {
	errors.TodoP4("improve performance by reusing slices")

	z.Etiketten = z1.Etiketten.ImmutableClone()

	z.Sorted = make([]string, len(z1.Sorted))
	copy(z.Sorted, z1.Sorted)
}
