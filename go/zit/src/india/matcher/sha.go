package matcher

import (
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type Sha kennung.Sha

func (t Sha) Each(_ schnittstellen.FuncIter[Matcher]) error {
	return nil
}

func (t Sha) MatcherLen() int {
	return 0
}

func (t Sha) ContainsMatchableExactly(m *sku.Transacted) bool {
	return t.ContainsMatchable(m)
}

func (t Sha) ContainsMatchable(m *sku.Transacted) bool {
	if kennung.Sha(t).EqualsSha(m.GetAkteSha()) {
		return true
	}

	return false
}
