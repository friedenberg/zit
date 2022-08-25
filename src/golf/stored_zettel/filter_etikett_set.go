package stored_zettel

import (
	"github.com/friedenberg/zit/src/delta/etikett"
)

type FilterEtikettSet struct {
	Or bool
	etikett.Set
}

func (f FilterEtikettSet) IncludeNamedZettel(z Named) bool {
	set := z.Stored.Zettel.Etiketten.IntersectPrefixes(f.Set)

	if f.Or {
		//at least one of the etiketten matches, resolving to a true or
		return set.Len() > 0
	} else {
		//by checking equal or greater than, we include zettels that have multiple
		//matches to the original set
		return set.Len() >= f.Set.Len()
	}
}
