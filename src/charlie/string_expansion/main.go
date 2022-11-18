package string_expansion

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

type Expander[T collections.ProtoObjekte, T1 interface {
	*T
	collections.ProtoObjektePointer
}] interface {
	Expand(string) collections.ValueSet[T, T1]
}
