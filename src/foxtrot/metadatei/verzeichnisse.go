package metadatei

import (
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type Verzeichnisse struct {
	Archiviert        values.Bool
	ExpandedEtiketten kennung.EtikettSet
	ImplicitEtiketten kennung.EtikettSet
	Mutter            [2]sha.Sha // shas of parent Metadatei
	Sha               sha.Sha    // sha of Metadatei
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
	v.Mutter[0].Reset()
	v.Mutter[1].Reset()
	v.Sha.Reset()
}

func (a *Verzeichnisse) ResetWith(b *Verzeichnisse) {
	a.Archiviert.ResetWith(b.Archiviert)
	a.ImplicitEtiketten = b.ImplicitEtiketten
	a.ExpandedEtiketten = b.ExpandedEtiketten
	a.Mutter = b.Mutter
	a.Sha = b.Sha
}
