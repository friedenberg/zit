package matcher

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type FD kennung.FD

func (fd FD) String() string {
	return kennung.FD(fd).String()
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

func (fd FD) ContainsMatchable(m *sku.Transacted) (ok bool) {
	il := m.GetKennungLike()

	switch il.GetGattung() {
	case gattung.Zettel:
		var h kennung.Hinweis

		if h, ok = kennung.FD(fd).AsHinweis(); !ok {
			return false
		}

		ok := kennung.Equals(h, il)

		return ok

	default:
		errors.TodoP1("support other gattung")
	}

	return false
}
