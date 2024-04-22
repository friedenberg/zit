package query

import (
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type QueryWithHidden struct {
	Query
	Hidden sku.Query
}

func (q *QueryWithHidden) Clone() *QueryWithHidden {
	q1 := q.Query.Clone()

	return &QueryWithHidden{
		Query:  *q1,
		Hidden: q.Hidden,
	}
}

func (q *QueryWithHidden) ContainsSku(sk *sku.Transacted) bool {
	if q.ShouldHide(sk) {
		return false
	}

	return q.Query.ContainsSku(sk)
}

func (q *QueryWithHidden) ShouldHide(sk *sku.Transacted) bool {
	// this gets checked more than once for every sku, maybe merge querywithhidden
	// and query?
	_, ok := q.Kennung[sk.Kennung.String()]

	if q.IncludesHidden() || q.Hidden == nil || ok {
		return false
	}

	return q.Hidden.ContainsSku(sk)
}
