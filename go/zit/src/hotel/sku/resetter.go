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

func (externalResetter) ResetWith(a *External, b *External) {
	errors.PanicIfError(a.Transacted.ObjectId.ResetWithIdLike(&b.Transacted.ObjectId))
	object_metadata.Resetter.ResetWith(&a.Transacted.Metadata, &b.Transacted.Metadata)
	TransactedResetter.ResetWith(b.GetSku(), a.GetSku())
	b.Transacted.Fields = b.Transacted.Fields[:0]
	b.Transacted.Fields = append(b.Transacted.Fields, a.Transacted.Fields...)
	b.ExternalType = a.ExternalType
	b.RepoId = a.RepoId
	b.State = a.State
	b.ExternalObjectId.ResetWith(&a.ExternalObjectId)
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
