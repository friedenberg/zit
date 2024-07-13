package store_verzeichnisse

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type flushQueryGroup struct {
	kennung.Sigil
}

func (qg *flushQueryGroup) SetIncludeHistory() {
	qg.Add(kennung.SigilHistory)
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

func (s *flushQueryGroup) ContainsKennung(_ *kennung.Kennung2) bool {
	return false
}

func (s *flushQueryGroup) GetSigil() kennung.Sigil {
	return s.Sigil
}

func (s *flushQueryGroup) Each(_ interfaces.FuncIter[sku.Query]) error {
	return nil
}
