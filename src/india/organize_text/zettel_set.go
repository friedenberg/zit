package organize_text

import (
	"sort"

	"github.com/friedenberg/zit/src/bravo/collections"
)

func sortZettelSet(
	s collections.MutableValueSet[zettel, *zettel],
) (out []zettel) {
	out = s.Elements()

	sort.Slice(out, func(i, j int) bool {
		if out[i].Bezeichnung == out[j].Bezeichnung {
			return out[i].Hinweis.Less(out[j].Hinweis)
		} else {
			return out[i].Bezeichnung.Less(out[j].Bezeichnung)
		}
	})

	return
}
