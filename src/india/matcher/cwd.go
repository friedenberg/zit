package matcher

import (
	"github.com/friedenberg/zit/src/echo/fd"
)

type MatcherCwd interface {
	Matcher
	GetCwdFDs() fd.Set
}

type matcherCwdNop struct {
	Matcher
}

func (_ matcherCwdNop) GetCwdFDs() fd.Set {
	return fd.MakeSet()
}

func MakeMatcherCwdNop(m Matcher) MatcherCwd {
	return matcherCwdNop{Matcher: m}
}
