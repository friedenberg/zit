package matcher

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type archiviert struct{}

func MakeArchiviert() Matcher {
	return &archiviert{}
}

func (m archiviert) MatcherLen() int {
	return 0
}

func (m archiviert) String() string {
	return ""
}

func (matcher archiviert) ContainsMatchable(matchable *sku.Transacted) bool {
	if !matchable.GetMetadatei().Verzeichnisse.Archiviert.Bool() {
		return false
	}

	return true
}

func (matcher archiviert) Each(f schnittstellen.FuncIter[Matcher]) error {
	return nil
}
