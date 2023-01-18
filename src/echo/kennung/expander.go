package kennung

import (
	"regexp"

	"github.com/friedenberg/zit/src/delta/collections"
)

var (
	regexExpandTagsHyphens *regexp.Regexp
	ExpanderRight          Expander
	ExpanderAll            Expander
)

type Expander interface {
	Expand(collections.FuncSetString, string)
}

func init() {
	regexExpandTagsHyphens = regexp.MustCompile(`-`)
	ExpanderRight = MakeExpanderRight(`-`)
	ExpanderAll = MakeExpanderAll(`-`)
}
