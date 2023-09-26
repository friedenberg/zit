package matcher

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
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
