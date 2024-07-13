package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
)

var TransactedResetter transactedResetter

type transactedResetter struct{}

func (transactedResetter) Reset(a *Transacted) {
	a.Kopf.Reset()
	a.Kennung.SetGenre(genres.Unknown)
	object_metadata.Resetter.Reset(&a.Metadatei)
	a.TransactionIndex.Reset()
}

func (transactedResetter) ResetWith(a *Transacted, b *Transacted) {
	a.Kopf = b.Kopf
	errors.PanicIfError(a.Kennung.ResetWithIdLike(&b.Kennung))
	object_metadata.Resetter.ResetWith(&a.Metadatei, &b.Metadatei)
	a.TransactionIndex.SetInt(b.TransactionIndex.Int())
}
