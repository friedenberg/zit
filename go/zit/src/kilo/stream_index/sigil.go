package stream_index

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type flushQueryGroup struct {
	ids.Sigil
}

func (qg *flushQueryGroup) SetIncludeHistory() {
	qg.Add(ids.SigilHistory)
}

func (qg *flushQueryGroup) HasHidden() bool {
	return false
}

func (qg *flushQueryGroup) Get(_ genres.Genre) (sku.QueryWithSigilAndObjectId, bool) {
	return qg, true
}

func (s *flushQueryGroup) ContainsSku(_ *sku.Transacted) bool {
	return true
}

func (s *flushQueryGroup) ContainsObjectId(_ *ids.ObjectId) bool {
	return false
}

func (s *flushQueryGroup) GetSigil() ids.Sigil {
	return s.Sigil
}

func (s *flushQueryGroup) Each(_ interfaces.FuncIter[sku.Query]) error {
	return nil
}
