package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
)

var TransactedResetter transactedResetter

type transactedResetter struct{}

func (transactedResetter) Reset(a *Transacted) {
	a.ObjectId.SetGenre(genres.Unknown)
	object_metadata.Resetter.Reset(&a.Metadata)
}

func (transactedResetter) ResetWith(a *Transacted, b *Transacted) {
	errors.PanicIfError(a.ObjectId.ResetWithIdLike(&b.ObjectId))
	object_metadata.Resetter.ResetWith(&a.Metadata, &b.Metadata)
}

var ExternalResetter externalResetter

type externalResetter struct{}

func (externalResetter) Reset(a *External) {
	a.ObjectId.SetGenre(genres.Unknown)
	object_metadata.Resetter.Reset(&a.Metadata)
}

func (externalResetter) ResetWith(a *External, b *External) {
	errors.PanicIfError(a.ObjectId.ResetWithIdLike(&b.ObjectId))
	object_metadata.Resetter.ResetWith(&a.Metadata, &b.Metadata)
}
