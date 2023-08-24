package matcher

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type MatcherCwd interface {
	Matcher
	GetCwdFDs() schnittstellen.SetLike[kennung.FD]
}

type matcherCwdNop struct {
	Matcher
}

func (_ matcherCwdNop) GetCwdFDs() schnittstellen.SetLike[kennung.FD] {
	return collections_value.MakeValueSet[kennung.FD](nil)
}

func MakeMatcherCwdNop(m Matcher) MatcherCwd {
	return matcherCwdNop{Matcher: m}
}
