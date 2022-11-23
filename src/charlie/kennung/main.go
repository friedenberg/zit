package kennung

import (
	"regexp"

	"github.com/friedenberg/zit/src/bravo/collections"
)

var (
	regexExpandTagsHyphens *regexp.Regexp
)

func init() {
	regexExpandTagsHyphens = regexp.MustCompile(`-`)
}

type Expander[T collections.ValueElement, T1 collections.ValueElementPtr[T]] interface {
	Expand(string) collections.ValueSet[T, T1]
}
