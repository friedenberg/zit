package object_metadata

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/tag_paths"
)

type Verzeichnisse struct {
	Schlummernd       values.Bool
	ExpandedEtiketten ids.TagMutableSet // public for gob, but should be private
	ImplicitEtiketten ids.TagMutableSet // public for gob, but should be private
	Etiketten         tag_paths.Etiketten
	QueryPath
}

func (v *Verzeichnisse) GetExpandedEtiketten() ids.TagSet {
	return v.GetExpandedEtikettenMutable()
}

func (v *Verzeichnisse) AddEtikettExpandedPtr(e *ids.Tag) (err error) {
	return iter.AddClonePool(
		v.GetExpandedEtikettenMutable(),
		ids.GetTagPool(),
		ids.TagResetter,
		e,
	)
}

func (v *Verzeichnisse) GetExpandedEtikettenMutable() ids.TagMutableSet {
	if v.ExpandedEtiketten == nil {
		v.ExpandedEtiketten = ids.MakeTagMutableSet()
	}

	return v.ExpandedEtiketten
}

func (v *Verzeichnisse) SetExpandedEtiketten(e ids.TagSet) {
	es := v.GetExpandedEtikettenMutable()
	iter.ResetMutableSetWithPool(es, ids.GetTagPool())

	if e == nil {
		return
	}

	errors.PanicIfError(e.Each(es.Add))
}

func (v *Verzeichnisse) GetImplicitEtiketten() ids.TagSet {
	return v.GetImplicitEtikettenMutable()
}

func (v *Verzeichnisse) AddEtikettImplicitPtr(e *ids.Tag) (err error) {
	return iter.AddClonePool(
		v.GetImplicitEtikettenMutable(),
		ids.GetTagPool(),
		ids.TagResetter,
		e,
	)
}

func (v *Verzeichnisse) GetImplicitEtikettenMutable() ids.TagMutableSet {
	if v.ImplicitEtiketten == nil {
		v.ImplicitEtiketten = ids.MakeTagMutableSet()
	}

	return v.ImplicitEtiketten
}

func (v *Verzeichnisse) SetImplicitEtiketten(e ids.TagSet) {
	es := v.GetImplicitEtikettenMutable()
	iter.ResetMutableSetWithPool(es, ids.GetTagPool())

	if e == nil {
		return
	}

	errors.PanicIfError(e.Each(es.Add))
}
