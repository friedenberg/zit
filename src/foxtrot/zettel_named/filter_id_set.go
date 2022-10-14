package zettel_named

import (
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/delta/id_set"
)

type FilterIdSet struct {
	id_set.Set
	Or bool
}

//TODO improve the performance of this query
func (f FilterIdSet) IncludeNamedZettel(z *Zettel) (ok bool) {
	needsEt := f.Set.Etiketten().Len() > 0
	okEt := false

	expanded := z.Stored.Zettel.Etiketten.Expanded(etikett.ExpanderRight{})

LOOP:
	for _, e := range f.Set.Etiketten().Sorted() {
		okEt = expanded.Contains(e)

		switch {
		case !okEt && !f.Or:
			break LOOP

		case okEt && f.Or:
			break LOOP

		default:
			continue
		}
	}

	needsTyp := len(f.Set.Typen()) > 0
	okTyp := false

	for _, t := range f.Set.Typen() {
		if okTyp = t.Includes(z.Stored.Zettel.Typ.Etikett); okTyp {
			break
		}
	}

	needsHin := len(f.Set.Hinweisen()) > 0
	okHin := false || len(f.Set.Hinweisen()) == 0

	for _, h := range f.Set.Hinweisen() {
		if okHin = h.Equals(z.Hinweis); okHin {
			break
		}
	}

	switch {
	case f.Or:
		ok = (okHin && needsHin) || (okTyp && needsTyp) || (okEt && needsEt)

	default:
		ok = (okHin || !needsHin) && (okTyp || !needsTyp) && (okEt || !needsEt)
	}

	return
}
