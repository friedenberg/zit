package matcher

import (
	"code.linenisgreat.com/zit/src/echo/fd"
)

type matcherCwd interface {
	Matcher
	GetCwdFDs() fd.Set
}

type matcherCwdNop struct {
	Matcher
}

func (_ matcherCwdNop) GetCwdFDs() fd.Set {
	return fd.MakeSet()
}

func MakeMatcherCwdNop(m Matcher) matcherCwd {
	return matcherCwdNop{Matcher: m}
}
