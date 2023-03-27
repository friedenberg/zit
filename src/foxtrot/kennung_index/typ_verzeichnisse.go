package kennung_index

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type TypVerzeichnisse struct {
	Expanded []string
}

func (z *TypVerzeichnisse) ResetWithTyp(t kennung.Typ) {
	ex := kennung.ExpandOne(t, kennung.ExpanderAll)
	z.Expanded = collections.SortedStrings[kennung.Typ](ex)
}

func (z *TypVerzeichnisse) Reset() {
	z.Expanded = []string{}
}

func (z *TypVerzeichnisse) ResetWith(z1 TypVerzeichnisse) {
	errors.TodoP4("improve performance by reusing slices")

	z.Expanded = make([]string, len(z1.Expanded))
	copy(z.Expanded, z1.Expanded)
}
