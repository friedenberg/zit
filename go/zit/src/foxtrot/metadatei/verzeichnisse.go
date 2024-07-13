package metadatei

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/etiketten_path"
)

type Verzeichnisse struct {
	Schlummernd       values.Bool
	ExpandedEtiketten kennung.TagMutableSet // public for gob, but should be private
	ImplicitEtiketten kennung.TagMutableSet // public for gob, but should be private
	Etiketten         etiketten_path.Etiketten
	QueryPath
}

func (v *Verzeichnisse) GetExpandedEtiketten() kennung.TagSet {
	return v.GetExpandedEtikettenMutable()
}

func (v *Verzeichnisse) AddEtikettExpandedPtr(e *kennung.Tag) (err error) {
	return iter.AddClonePool(
		v.GetExpandedEtikettenMutable(),
		kennung.GetTagPool(),
		kennung.TagResetter,
		e,
	)
}

func (v *Verzeichnisse) GetExpandedEtikettenMutable() kennung.TagMutableSet {
	if v.ExpandedEtiketten == nil {
		v.ExpandedEtiketten = kennung.MakeTagMutableSet()
	}

	return v.ExpandedEtiketten
}

func (v *Verzeichnisse) SetExpandedEtiketten(e kennung.TagSet) {
	es := v.GetExpandedEtikettenMutable()
	iter.ResetMutableSetWithPool(es, kennung.GetTagPool())

	if e == nil {
		return
	}

	errors.PanicIfError(e.Each(es.Add))
}

func (v *Verzeichnisse) GetImplicitEtiketten() kennung.TagSet {
	return v.GetImplicitEtikettenMutable()
}

func (v *Verzeichnisse) AddEtikettImplicitPtr(e *kennung.Tag) (err error) {
	return iter.AddClonePool(
		v.GetImplicitEtikettenMutable(),
		kennung.GetTagPool(),
		kennung.TagResetter,
		e,
	)
}

func (v *Verzeichnisse) GetImplicitEtikettenMutable() kennung.TagMutableSet {
	if v.ImplicitEtiketten == nil {
		v.ImplicitEtiketten = kennung.MakeTagMutableSet()
	}

	return v.ImplicitEtiketten
}

func (v *Verzeichnisse) SetImplicitEtiketten(e kennung.TagSet) {
	es := v.GetImplicitEtikettenMutable()
	iter.ResetMutableSetWithPool(es, kennung.GetTagPool())

	if e == nil {
		return
	}

	errors.PanicIfError(e.Each(es.Add))
}
