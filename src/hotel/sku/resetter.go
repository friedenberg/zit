package sku

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

var TransactedResetter transactedResetter

type transactedResetter struct{}

func (transactedResetter) Reset(a *Transacted) {
	a.Kopf.Reset()
	a.ObjekteSha.Reset()
	a.Kennung.SetGattung(gattung.Unknown)
	metadatei.Resetter.Reset(&a.Metadatei)
	a.TransactionIndex.Reset()
}

func (transactedResetter) ResetWithPtr(a *Transacted, b *Transacted) {
	a.Kopf = b.Kopf
	errors.PanicIfError(a.ObjekteSha.SetShaLike(&b.ObjekteSha))
	errors.PanicIfError(a.Kennung.ResetWithKennung(&b.Kennung))
	metadatei.Resetter.ResetWithPtr(&a.Metadatei, &b.Metadatei)
	a.TransactionIndex.SetInt(b.TransactionIndex.Int())
}
