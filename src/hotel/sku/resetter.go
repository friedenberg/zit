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

func (r transactedResetter) ResetWith(a *Transacted, b Transacted) {
	a.Kopf = b.Kopf
	a.ObjekteSha = b.ObjekteSha
	a.Kennung.ResetWithKennung(b.Kennung)
	metadatei.Resetter.ResetWith(&a.Metadatei, b.Metadatei)
	a.TransactionIndex.SetInt(b.TransactionIndex.Int())
}

func (transactedResetter) ResetWithPtr(a *Transacted, b *Transacted) {
	a.Kopf = b.Kopf
	a.ObjekteSha = b.ObjekteSha
	errors.PanicIfError(a.Kennung.ResetWithKennung(&b.Kennung))
	metadatei.Resetter.ResetWithPtr(&a.Metadatei, &b.Metadatei)
	a.TransactionIndex.SetInt(b.TransactionIndex.Int())
}