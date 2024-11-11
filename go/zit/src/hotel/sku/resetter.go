package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
)

var TransactedResetter transactedResetter

type transactedResetter struct{}

func (transactedResetter) Reset(dst *Transacted) {
	dst.ObjectId.Reset()
	object_metadata.Resetter.Reset(&dst.Metadata)
	dst.ExternalType.Reset()
	dst.RepoId.Reset()
	dst.State = external_state.Unknown
}

func (transactedResetter) ResetWith(dst *Transacted, src *Transacted) {
	errors.PanicIfError(dst.ObjectId.ResetWithIdLike(&src.ObjectId))
	object_metadata.Resetter.ResetWith(&dst.Metadata, &src.Metadata)
	dst.ExternalType = src.ExternalType
	dst.RepoId = src.RepoId
	dst.State = src.State
	dst.ExternalObjectId.ResetWith(&src.ExternalObjectId)
}

func (transactedResetter) ResetWithExceptFields(dst *Transacted, src *Transacted) {
	errors.PanicIfError(dst.ObjectId.ResetWithIdLike(&src.ObjectId))
	object_metadata.Resetter.ResetWithExceptFields(&dst.Metadata, &src.Metadata)
	dst.ExternalType = src.ExternalType
	dst.RepoId = src.RepoId
	dst.State = src.State
	dst.ExternalObjectId.ResetWith(&src.ExternalObjectId)
}

var Resetter resetter

type resetter struct{}

func (resetter) Reset(sl TransactedGetter) {
	TransactedResetter.Reset(sl.GetSku())
}

func (resetter) ResetWith(dst, src TransactedGetter) {
	TransactedResetter.ResetWith(dst.GetSku(), src.GetSku())
}

var CheckedOutResetter checkedOutResetter

type checkedOutResetter struct{}

func (checkedOutResetter) Reset(dst *CheckedOut) {
	TransactedResetter.Reset(dst.GetSku())
	TransactedResetter.Reset(dst.GetSkuExternalLike().GetSku())
	dst.SetState(checked_out_state.Unknown)
	dst.IsImport = false
}

func (checkedOutResetter) ResetWith(dst *CheckedOut, src *CheckedOut) {
	TransactedResetter.ResetWith(dst.GetSku(), src.GetSku())
	TransactedResetter.ResetWith(dst.GetSkuExternalLike().GetSku(), src.GetSkuExternalLike().GetSku())
	dst.SetState(src.State)
	dst.IsImport = src.IsImport
}
