package matcher

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type FD struct {
	*fd.FD
}

func (f FD) String() string {
	return f.FD.String()
}

func (_ FD) Each(_ schnittstellen.FuncIter[Matcher]) error {
	return nil
}

func (fd FD) MatcherLen() int {
	return 0
}

func (fd FD) ContainsMatchableExactly(m *sku.Transacted) (ok bool) {
	return fd.ContainsMatchable(m)
}

func (f FD) ContainsMatchable(m *sku.Transacted) (ok bool) {
	il := m.GetKennung()

	switch il.GetGattung() {
	case gattung.Zettel:
		var h kennung.Hinweis

		if h, ok = kennung.AsHinweis(f.FD); !ok {
			return false
		}

		ok := kennung.Equals(h, il)

		return ok

	default:
		errors.TodoP1("support other gattung")
	}

	return false
}
