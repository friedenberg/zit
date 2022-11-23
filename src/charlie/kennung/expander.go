package kennung

import (
	"regexp"

	"github.com/friedenberg/zit/src/bravo/collections"
)

type ExpanderEtikett = Expander[Etikett, *Etikett]

var (
	regexExpandTagsHyphens *regexp.Regexp
	ExpanderEtikettRight   ExpanderEtikett
	ExpanderEtikettAll     ExpanderEtikett
)

type Expander[T collections.ValueElement, T1 collections.ValueElementPtr[T]] interface {
	Expand(string) collections.ValueSet[T, T1]
}

func init() {
	regexExpandTagsHyphens = regexp.MustCompile(`-`)
	ExpanderEtikettRight = MakeExpanderRight[Etikett, *Etikett](`-`)
	ExpanderEtikettAll = MakeExpanderAll[Etikett, *Etikett](`-`)
}
