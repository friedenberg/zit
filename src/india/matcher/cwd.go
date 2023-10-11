package matcher

import (
	"github.com/friedenberg/zit/src/echo/kennung"
)

type MatcherCwd interface {
	Matcher
	GetCwdFDs() kennung.FDSet
}

type matcherCwdNop struct {
	Matcher
}

func (_ matcherCwdNop) GetCwdFDs() kennung.FDSet {
	return kennung.MakeFDSet()
}

func MakeMatcherCwdNop(m Matcher) MatcherCwd {
	return matcherCwdNop{Matcher: m}
}
