package matcher

import (
	"sync/atomic"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type Archiviert interface {
	Matcher
	Count() int64
}

type archiviert struct {
	count int64
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

func (matcher *archiviert) ContainsMatchable(matchable *sku.Transacted) bool {
	if !matchable.GetMetadatei().Verzeichnisse.Archiviert.Bool() {
		atomic.AddInt64(&matcher.count, 1)
		return false
	}

	return true
}

func (matcher archiviert) Each(f schnittstellen.FuncIter[Matcher]) error {
	return nil
}
