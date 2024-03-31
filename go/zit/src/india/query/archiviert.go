package query

import (
	"sync/atomic"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type Archiviert interface {
	sku.Query
	CountArchiviert() int64
	Count() int64
}

type archiviert struct {
	count, countArchiviert int64
}

func MakeArchiviert() Archiviert {
	return &archiviert{}
}

func (m archiviert) MatcherLen() int {
	return 0
}

func (m archiviert) String() string {
	return ""
}

func (m archiviert) Count() int64 {
	return m.count
}

func (m archiviert) CountArchiviert() int64 {
	return m.countArchiviert
}

func (matcher *archiviert) ContainsSku(matchable *sku.Transacted) bool {
	if !matchable.GetMetadatei().Verzeichnisse.Archiviert.Bool() {
		atomic.AddInt64(&matcher.countArchiviert, 1)
		return false
	}

	atomic.AddInt64(&matcher.count, 1)

	return true
}

func (matcher archiviert) Each(f schnittstellen.FuncIter[sku.Query]) error {
	return nil
}
