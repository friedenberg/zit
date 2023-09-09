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

func (v *Verzeichnisse) GetExpandedEtiketten() kennung.EtikettSet {
	if v.ExpandedEtiketten == nil {
		return kennung.EtikettSetEmpty
	}

	return v.ExpandedEtiketten
}

func (v *Verzeichnisse) GetImplicitEtiketten() kennung.EtikettSet {
	if v.ImplicitEtiketten == nil {
		return kennung.EtikettSetEmpty
	}

	return v.ImplicitEtiketten
}

func (v *Verzeichnisse) Reset() {
	v.Archiviert.Reset()
	v.ImplicitEtiketten = kennung.EtikettSetEmpty
	v.ExpandedEtiketten = kennung.EtikettSetEmpty
}

func (a *Verzeichnisse) ResetWith(b *Verzeichnisse) {
	a.Archiviert.ResetWith(b.Archiviert)
	a.ImplicitEtiketten = b.ImplicitEtiketten
	a.ExpandedEtiketten = b.ExpandedEtiketten
}
