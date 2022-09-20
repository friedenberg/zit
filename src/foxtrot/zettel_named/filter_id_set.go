package zettel_named

import (
	"github.com/friedenberg/zit/src/delta/id_set"
)

type FilterIdSet struct {
	id_set.Set
	And bool
}

func (f FilterIdSet) IncludeNamedZettel(z Zettel) (ok bool) {
	for _, t := range f.Set.Typen() {
		ok = t.Includes(z.Stored.Zettel.Typ.Etikett)

		switch {
		case ok && f.And:
			continue

		case f.And:
		case ok:
			return
		}
	}

	for _, h := range f.Set.Hinweisen() {
		ok = h.Equals(z.Hinweis)

		switch {
		case ok && f.And:
			continue

		case f.And:
		case ok:
			return
		}
	}

	for _, e := range f.Set.Etiketten().Sorted() {
		ok = z.Stored.Zettel.Etiketten.Contains(e)

		switch {
		case ok && f.And:
			continue

		case f.And:
		case ok:
			return
		}
	}

	return
}
