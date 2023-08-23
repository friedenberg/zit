package matcher

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type Sha kennung.Sha

func (t Sha) Each(_ schnittstellen.FuncIter[Matcher]) error {
	return nil
}

func (t Sha) MatcherLen() int {
	return 0
}

func (t Sha) ContainsMatchableExactly(m Matchable) bool {
	return t.ContainsMatchable(m)
}

func (t Sha) ContainsMatchable(m Matchable) bool {
	if kennung.Sha(t).EqualsSha(m.GetAkteSha()) {
		return true
	}

	return false
}
