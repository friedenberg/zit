package store_fs

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO migrate to *sku.External
type External struct {
	sku.External
	item Item // TODO remove and keep separately
}

func (t *External) GetSkuExternalLike() sku.ExternalLike {
	return t
}

func (a *External) Clone() sku.ExternalLike {
	b := GetExternalPool().Get()
	sku.ExternalResetter.ResetWith(&b.External, &a.External)
	b.item.ResetWith(&a.item)
	return b
}

func (c *External) GetSku() *sku.Transacted {
	return &c.External.Transacted.Transacted
}

func (a *External) ResetWith(b *External) {
	a.item.ResetWith(b.GetFDs())
	sku.Resetter.ResetWith(a, b)
}

func (a *External) String() string {
	return fmt.Sprintf(
		". %s %s %s %s",
		a.GetGenre(),
		a.GetObjectId(),
		a.Transacted.GetObjectSha(),
		a.GetBlobSha(),
	)
}

func (a *External) SetBlobSha(v interfaces.Sha) (err error) {
	if err = a.Transacted.Metadata.Blob.SetShaLike(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = a.item.Blob.SetShaLike(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *External) GetFDs() *Item {
	return &a.item
}

func (a *External) GetBlobFD() *fd.FD {
	return &a.item.Blob
}

func (a *External) SetBlobFD(v *fd.FD) {
	a.item.Blob.ResetWith(v)
	a.Transacted.Metadata.Blob.SetShaLike(v.GetShaLike())
}

func (a *External) GetBlobPath() string {
	return a.item.Blob.GetPath()
}

func (o *External) GetKey() string {
	return fmt.Sprintf("%s.%s", o.GetGenre(), o.GetObjectId())
}

func GetCheckoutModeOrError(
	el sku.ExternalLike,
	originalMode checkout_mode.Mode,
) (m checkout_mode.Mode, err error) {
	var fds Item

	if err = fds.ReadFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	switch {
	case !fds.Object.IsEmpty() && !fds.Blob.IsEmpty():
		m = checkout_mode.MetadataAndBlob

	case !fds.Blob.IsEmpty():
		m = checkout_mode.BlobOnly

	case !fds.Object.IsEmpty():
		m = checkout_mode.MetadataOnly

	default:
		if fds.State == external_state.Recognized {
			m = checkout_mode.BlobRecognized
			return
		}

		err = checkout_mode.MakeErrInvalidCheckoutMode(
			errors.Errorf("all FD's are empty: %s", fds.Debug()),
		)
	}

	return
}

func GetConflictOrError(
	el sku.ExternalLike,
) (f *fd.FD, err error) {
	var fds Item

	if err = fds.ReadFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	f = &fds.Conflict

	return
}

func GetObjectOrError(
	el sku.ExternalLike,
) (f *fd.FD, err error) {
	var fds Item

	if err = fds.ReadFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	f = &fds.Object

	return
}

func SetObjectOrError(
	el sku.ExternalLike,
	object *fd.FD,
) (err error) {
	var fds Item

	if err = fds.ReadFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	fds.Object.ResetWith(object)

	if err = fds.WriteToExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func GetBlobOrError(
	el sku.ExternalLike,
) (f *fd.FD, err error) {
	var fds Item

	if err = fds.ReadFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	f = &fds.Blob

	return
}

func SetBlobOrError(
	el sku.ExternalLike,
	blob *fd.FD,
) (err error) {
	var fds Item

	if err = fds.ReadFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	fds.Blob.ResetWith(blob)

	if err = fds.WriteToExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func UpdateDescriptionFromBlobs(
	el sku.ExternalLike,
) (err error) {
	var fds Item

	if err = fds.ReadFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	sorted := iter.ElementsSorted(
		fds.MutableSetLike,
		func(a, b *fd.FD) bool {
			return a.GetPath() < b.GetPath()
		},
	)

	for _, f := range sorted {
		desc := f.FileNameSansExt()

		if err = el.GetSku().Metadata.Description.Set(desc); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = fds.WriteToExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
