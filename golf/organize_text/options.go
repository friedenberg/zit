package organize_text

import (
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

type Grouper interface {
	GroupZettel(stored_zettel.Named) []etikett.Set
}

type Sorter interface {
	SortGroups(etikett.Set, etikett.Set) bool
	SortZettels(stored_zettel.Named, stored_zettel.Named) bool
}

type Options struct {
	RootEtiketten etikett.Set
	GroupingEtiketten etikett.Slice
	// TODO option for combining or separating roots
	// Grouper
	// Sorter
}
