package kennung

import (
	"regexp"
)

var (
	regexExpandTagsHyphens *regexp.Regexp
	ExpanderRight          Expander
	ExpanderAll            Expander
)

type stringAdder interface {
	AddString(string) error
}

type Expander interface {
	Expand(stringAdder, string)
}

func init() {
	regexExpandTagsHyphens = regexp.MustCompile(`-`)
	ExpanderRight = MakeExpanderRight(`-`)
	ExpanderAll = MakeExpanderAll(`-`)
}
