package metadatei

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/values"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/etiketten_path"
)

type Verzeichnisse struct {
	Schlummernd        values.Bool
	ExpandedEtiketten kennung.EtikettMutableSet // public for gob, but should be private
	ImplicitEtiketten kennung.EtikettMutableSet // public for gob, but should be private
	Etiketten         etiketten_path.Etiketten
	QueryPath
}

func (v *Verzeichnisse) GetExpandedEtiketten() kennung.EtikettSet {
	return v.GetExpandedEtikettenMutable()
}

func (v *Verzeichnisse) AddEtikettExpandedPtr(e *kennung.Etikett) (err error) {
	return iter.AddClonePool(
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
	return iter.AddClonePool(
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
