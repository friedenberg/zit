package matcher

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
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
