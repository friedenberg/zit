package etiketten_path

import (
	"slices"

	"code.linenisgreat.com/zit/src/delta/catgut"
)

type (
	Etikett            = catgut.String
	EtikettWithParents struct {
		*Etikett
		Parents []*Path
	}
)

func (ewp *EtikettWithParents) AddParent(p *Path) {
	if p == nil || p.Len() == 1 {
		return
	}

	idxPath, okPath := slices.BinarySearchFunc(
		ewp.Parents,
		p,
		func(ep *Path, el *Path) int {
			return ep.Compare(p)
		},
	)

	if !okPath {
		ewp.Parents = slices.Insert(ewp.Parents, idxPath, p)
	}
}
