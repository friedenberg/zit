package store_fs

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type External struct {
	sku.Transacted
	FDs FDSet
}

func (t *External) GetSkuExternalLike() sku.ExternalLike {
	return t
}

func (t *External) GetExternalState() external_state.State {
	return t.FDs.State
}

func (a *External) Clone() sku.ExternalLike {
	b := GetExternalPool().Get()
	sku.TransactedResetter.ResetWith(&b.Transacted, &a.Transacted)
	b.FDs.ResetWith(&a.FDs)
	return b
}

func (c *External) GetSku() *sku.Transacted {
	return &c.Transacted
}

func (a *External) ResetWith(b *External) {
	a.FDs.ResetWith(b.GetFDs())
	sku.Resetter.ResetWith(a, b)
}

func (a *External) GetObjectId() *ids.ObjectId {
	return &a.ObjectId
}

func (a *External) GetMetadata() *object_metadata.Metadata {
	return &a.Metadata
}

func (a *External) GetGattung() interfaces.Genre {
	return a.ObjectId.GetGenre()
}

func (a *External) String() string {
	return fmt.Sprintf(
		". %s %s %s %s",
		a.GetGattung(),
		a.GetObjectId(),
		a.GetObjectSha(),
		a.GetBlobSha(),
	)
}

func (a *External) GetBlobSha() interfaces.Sha {
	return &a.Metadata.Blob
}

func (a *External) SetBlobSha(v interfaces.Sha) (err error) {
	if err = a.Metadata.Blob.SetShaLike(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = a.FDs.Blob.SetShaLike(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *External) GetFDs() *FDSet {
	return &a.FDs
}

func (a *External) GetBlobFD() *fd.FD {
	return &a.FDs.Blob
}

func (a *External) SetBlobFD(v *fd.FD) {
	a.FDs.Blob.ResetWith(v)
	a.Metadata.Blob.SetShaLike(v.GetShaLike())
}

func (a *External) GetBlobPath() string {
	return a.FDs.Blob.GetPath()
}

func (a *External) GetObjectFD() *fd.FD {
	return &a.FDs.Object
}

func (a *External) ResetWithExternalMaybe(
	b *FDSet,
) (err error) {
	k := &b.ObjectId
	a.ObjectId.ResetWithIdLike(k)
	object_metadata.Resetter.Reset(&a.Metadata)
	a.FDs.ResetWith(b)
	return
}

func (o *External) GetKey() string {
	return fmt.Sprintf("%s.%s", o.GetGattung(), o.GetObjectId())
}

func (e *External) GetCheckoutMode() (m checkout_mode.Mode, err error) {
	switch {
	case !e.FDs.Object.IsEmpty() && !e.FDs.Blob.IsEmpty():
		m = checkout_mode.MetadataAndBlob

	case !e.FDs.Blob.IsEmpty():
		m = checkout_mode.BlobOnly

	case !e.FDs.Object.IsEmpty():
		m = checkout_mode.MetadataOnly

	default:
		err = checkout_mode.MakeErrInvalidCheckoutMode(
			errors.Errorf("all FD's are empty"),
		)
	}

	return
}

type lessorExternal struct{}

func (lessorExternal) Less(a, b External) bool {
	panic("not supported")
}

func (lessorExternal) LessPtr(a, b *External) bool {
	return a.GetTai().Less(b.GetTai())
}

type equalerExternal struct{}

func (equalerExternal) Equals(a, b External) bool {
	panic("not supported")
}
