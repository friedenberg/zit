package metadatei

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type Verzeichnisse struct {
	Archiviert        values.Bool
	ExpandedEtiketten kennung.EtikettMutableSet // public for gob, but should be private
	ImplicitEtiketten kennung.EtikettMutableSet // public for gob, but should be private
	Mutter            sha.Sha                   // sha of parent Metadatei
}

func (v *Verzeichnisse) GetExpandedEtiketten() kennung.EtikettSet {
	return v.GetExpandedEtikettenMutable()
}

func (v *Verzeichnisse) AddEtikettExpandedPtr(e *kennung.Etikett) (err error) {
	return iter.AddClonePool[kennung.Etikett, *kennung.Etikett](
		v.GetExpandedEtikettenMutable(),
		kennung.GetEtikettPool(),
		kennung.EtikettResetter,
		e,
	)
}

func (v *Verzeichnisse) GetExpandedEtikettenMutable() kennung.EtikettMutableSet {
	if v.ExpandedEtiketten == nil {
		v.ExpandedEtiketten = kennung.MakeEtikettMutableSet()
	}

	return v.ExpandedEtiketten
}

func (v *Verzeichnisse) SetExpandedEtiketten(e kennung.EtikettSet) {
	es := v.GetExpandedEtikettenMutable()
	iter.ResetMutableSetWithPool(es, kennung.GetEtikettPool())

	if e == nil {
		return
	}

	errors.PanicIfError(e.Each(es.Add))
}

func (v *Verzeichnisse) GetImplicitEtiketten() kennung.EtikettSet {
	return v.GetImplicitEtikettenMutable()
}

func (v *Verzeichnisse) AddEtikettImplicitPtr(e *kennung.Etikett) (err error) {
	return iter.AddClonePool[kennung.Etikett, *kennung.Etikett](
		v.GetImplicitEtikettenMutable(),
		kennung.GetEtikettPool(),
		kennung.EtikettResetter,
		e,
	)
}

func (v *Verzeichnisse) GetImplicitEtikettenMutable() kennung.EtikettMutableSet {
	if v.ImplicitEtiketten == nil {
		v.ImplicitEtiketten = kennung.MakeEtikettMutableSet()
	}

	return v.ImplicitEtiketten
}

func (v *Verzeichnisse) SetImplicitEtiketten(e kennung.EtikettSet) {
	es := v.GetImplicitEtikettenMutable()
	iter.ResetMutableSetWithPool(es, kennung.GetEtikettPool())

	if e == nil {
		return
	}

	errors.PanicIfError(e.Each(es.Add))
}
