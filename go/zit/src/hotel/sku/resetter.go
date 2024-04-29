package sku

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/checked_out_state"
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

var CheckedOutResetter checkedOutResetter

type checkedOutResetter struct{}

func (checkedOutResetter) Reset(a *CheckedOut) {
	a.State = checked_out_state.StateUnknown
	a.IsImport = false
	a.Error = nil

	TransactedResetter.Reset(&a.Internal)
	TransactedResetter.Reset(&a.External.Transacted)
	a.External.FDs.Objekte.Reset()
	a.External.FDs.Akte.Reset()
}

func (checkedOutResetter) ResetWith(a *CheckedOut, b *CheckedOut) {
	a.State = b.State
	a.IsImport = b.IsImport
	a.Error = b.Error

	TransactedResetter.ResetWith(&a.Internal, &b.Internal)
	TransactedResetter.ResetWith(&a.External.Transacted, &b.External.Transacted)
	a.External.FDs.ResetWith(&b.External.FDs)
}
