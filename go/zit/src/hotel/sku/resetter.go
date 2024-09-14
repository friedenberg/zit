package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
)

var TransactedResetter transactedResetter

type transactedResetter struct{}

func (transactedResetter) Reset(a *Transacted) {
	a.ObjectId.Reset()
	object_metadata.Resetter.Reset(&a.Metadata)
}

func (transactedResetter) ResetWith(a *Transacted, b *Transacted) {
	errors.PanicIfError(a.ObjectId.ResetWithIdLike(&b.ObjectId))
	object_metadata.Resetter.ResetWith(&a.Metadata, &b.Metadata)
}

var ExternalResetter externalResetter

type externalResetter struct{}

func (externalResetter) Reset(a *External) {
	TransactedResetter.Reset(a.GetSku())
	a.Transacted.Fields = a.Transacted.Fields[:0]
	a.ExternalType.Reset()
	a.RepoId.Reset()
	a.State = external_state.Unknown
	a.ExternalObjectId.Reset()
}

func (externalResetter) ResetWith(dst *External, src *External) {
	errors.PanicIfError(dst.Transacted.ObjectId.ResetWithIdLike(&src.Transacted.ObjectId))
	object_metadata.Resetter.ResetWith(&dst.Transacted.Metadata, &src.Transacted.Metadata)
	TransactedResetter.ResetWith(dst.GetSku(), src.GetSku())
	dst.Transacted.Fields = dst.Transacted.Fields[:0]
	dst.Transacted.Fields = append(dst.Transacted.Fields, src.Transacted.Fields...)
	dst.ExternalType = src.ExternalType
	dst.RepoId = src.RepoId
	dst.State = src.State
	dst.ExternalObjectId.ResetWith(&src.ExternalObjectId)
}

var Resetter resetter

type resetter struct{}

func (resetter) Reset(sl TransactedGetter) {
	a := sl.GetSku()
	a.ObjectId.SetGenre(genres.Unknown)
	object_metadata.Resetter.Reset(&a.Metadata)
}

func (resetter) ResetWith(asl, bsl TransactedGetter) {
	a, b := asl.GetSku(), bsl.GetSku()
	errors.PanicIfError(a.ObjectId.ResetWithIdLike(&b.ObjectId))
	object_metadata.Resetter.ResetWith(&a.Metadata, &b.Metadata)
}
