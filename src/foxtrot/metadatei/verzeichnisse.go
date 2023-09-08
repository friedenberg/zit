package metadatei

import (
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type Verzeichnisse struct {
	Archiviert        values.Bool
	ExpandedEtiketten kennung.EtikettSet
	ImplicitEtiketten kennung.EtikettSet
}

func (v *Verzeichnisse) Reset() {
	v.Archiviert.Reset()
	v.ImplicitEtiketten = kennung.MakeEtikettSet()
	v.ExpandedEtiketten = kennung.MakeEtikettSet()
}

func (a *Verzeichnisse) ResetWith(b *Verzeichnisse) {
	a.Archiviert.ResetWith(b.Archiviert)
	a.ImplicitEtiketten = b.ImplicitEtiketten
	a.ExpandedEtiketten = b.ExpandedEtiketten
}
