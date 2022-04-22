package etikett

import "regexp"

var (
	regexExpandTagsHyphens *regexp.Regexp
)

func init() {
	regexExpandTagsHyphens = regexp.MustCompile(`-`)
}

type Expander interface {
	Expand(Etikett) Set
}
