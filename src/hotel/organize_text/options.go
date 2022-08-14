package organize_text

import (
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
)

type Grouper interface {
	GroupZettel(stored_zettel.Named) []etikett.Set
}

type Sorter interface {
	SortGroups(etikett.Set, etikett.Set) bool
	SortZettels(stored_zettel.Named, stored_zettel.Named) bool
}

type Options struct {
	RootEtiketten     etikett.Set
	GroupingEtiketten etikett.Slice
	ExtraEtiketten    etikett.Set
	// TODO option for combining or separating roots
	// Grouper
	// Sorter
}
