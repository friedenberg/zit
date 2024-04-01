package sku

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
)

var TransactedResetter transactedResetter

type transactedResetter struct{}

func (transactedResetter) Reset(a *Transacted) {
	a.Kopf.Reset()
	a.Kennung.SetGattung(gattung.Unknown)
	metadatei.Resetter.Reset(&a.Metadatei)
	a.TransactionIndex.Reset()
}

func (transactedResetter) ResetWith(a *Transacted, b *Transacted) {
	a.Kopf = b.Kopf
	errors.PanicIfError(a.Kennung.ResetWithKennung(&b.Kennung))
	metadatei.Resetter.ResetWith(&a.Metadatei, &b.Metadatei)
	a.TransactionIndex.SetInt(b.TransactionIndex.Int())
}
