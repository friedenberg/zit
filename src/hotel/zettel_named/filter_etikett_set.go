package zettel_named

import (
	"github.com/friedenberg/zit/src/delta/etikett"
)

type FilterEtikettSet struct {
	Or bool
	etikett.Set
}

func (f FilterEtikettSet) IncludeNamedZettel(z Zettel) (ok bool) {
	if f.Set.Len() == 0 {
		ok = true
		return
	}

	set := z.Stored.Zettel.Etiketten.IntersectPrefixes(f.Set)

	if f.Or {
		//at least one of the etiketten matches, resolving to a true or
		ok = set.Len() > 0
	} else {
		//by checking equal or greater than, we include zettels that have multiple
		//matches to the original set
		ok = set.Len() >= f.Set.Len()
	}

	return
}
