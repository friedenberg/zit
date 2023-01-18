package kennung

import (
	"regexp"

	"github.com/friedenberg/zit/src/schnittstellen"
)

var (
	regexExpandTagsHyphens *regexp.Regexp
	ExpanderRight          Expander
	ExpanderAll            Expander
)

type Expander interface {
	Expand(schnittstellen.FuncSetString, string)
}

func init() {
	regexExpandTagsHyphens = regexp.MustCompile(`-`)
	ExpanderRight = MakeExpanderRight(`-`)
	ExpanderAll = MakeExpanderAll(`-`)
}
