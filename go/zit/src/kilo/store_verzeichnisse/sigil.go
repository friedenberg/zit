package store_verzeichnisse

import (
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type sigil struct {
	kennung.Sigil
}

func (qg *sigil) Get(_ gattung.Gattung) (sku.QueryWithSigilAndKennung, bool) {
	return qg, true
}

func (s *sigil) ContainsSku(_ *sku.Transacted) bool {
	return true
}

func (s *sigil) String() string {
	panic("should never be called")
}

func (s *sigil) ContainsKennung(_ *kennung.Kennung2) bool {
	return false
}

func (s *sigil) GetSigil() kennung.Sigil {
	return s.Sigil
}

func (s *sigil) Each(_ schnittstellen.FuncIter[sku.Query]) error {
	return nil
}
