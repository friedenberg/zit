package store_verzeichnisse

import (
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type flushQueryGroup struct {
	kennung.Sigil
}

func (qg *flushQueryGroup) HasHidden() bool {
	return false
}

func (qg *flushQueryGroup) Get(_ gattung.Gattung) (sku.QueryWithSigilAndKennung, bool) {
	return qg, true
}

func (s *flushQueryGroup) ContainsSku(_ *sku.Transacted) bool {
	return true
}

func (s *flushQueryGroup) String() string {
	panic("should never be called")
}

func (s *flushQueryGroup) ContainsKennung(_ *kennung.Kennung2) bool {
	return false
}

func (s *flushQueryGroup) GetSigil() kennung.Sigil {
	return s.Sigil
}

func (s *flushQueryGroup) Each(_ schnittstellen.FuncIter[sku.Query]) error {
	return nil
}
