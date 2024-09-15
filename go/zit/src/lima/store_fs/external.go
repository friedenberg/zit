package store_fs

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO migrate to *sku.External
type External struct {
	sku.External
	fds Item // TODO remove and keep separately
}

func (t *External) GetSkuExternalLike() sku.ExternalLike {
	return t
}

func (t *External) GetExternalState() external_state.State {
	return t.fds.State
}

func (a *External) Clone() sku.ExternalLike {
	b := GetExternalPool().Get()
	sku.ExternalResetter.ResetWith(&b.External, &a.External)
	b.fds.ResetWith(&a.fds)
	return b
}

func (c *External) GetSku() *sku.Transacted {
	return &c.External.Transacted.Transacted
}

func (a *External) ResetWith(b *External) {
	a.fds.ResetWith(b.GetFDs())
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

	if err = a.fds.Blob.SetShaLike(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *External) GetFDs() *Item {
	return &a.fds
}

func (a *External) GetBlobFD() *fd.FD {
	return &a.fds.Blob
}

func (a *External) SetBlobFD(v *fd.FD) {
	a.fds.Blob.ResetWith(v)
	a.Transacted.Metadata.Blob.SetShaLike(v.GetShaLike())
}

func (a *External) GetBlobPath() string {
	return a.fds.Blob.GetPath()
}

func (a *External) GetObjectFD() *fd.FD {
	return &a.fds.Object
}

func (a *External) ResetWithExternalMaybe(
	b *Item,
) (err error) {
	k := &b.ObjectId
	a.Transacted.ObjectId.ResetWithIdLike(k)
	object_metadata.Resetter.Reset(&a.Transacted.Metadata)
	a.fds.ResetWith(b)
	return
}

func (o *External) GetKey() string {
	return fmt.Sprintf("%s.%s", o.GetGenre(), o.GetObjectId())
}

func GetCheckoutModeOrError(
	el sku.ExternalLike,
	originalMode checkout_mode.Mode,
) (m checkout_mode.Mode, err error) {
	var e *External
	ok := false

	if e, ok = el.(*External); !ok {
		m = originalMode
		return
	}

	switch {
	case !e.fds.Object.IsEmpty() && !e.fds.Blob.IsEmpty():
		m = checkout_mode.MetadataAndBlob

	case !e.fds.Blob.IsEmpty():
		m = checkout_mode.BlobOnly

	case !e.fds.Object.IsEmpty():
		m = checkout_mode.MetadataOnly

	default:
		if e.fds.State == external_state.Recognized {
			m = checkout_mode.BlobRecognized
			return
		}

		err = checkout_mode.MakeErrInvalidCheckoutMode(
			errors.Errorf("all FD's are empty: %s", e.fds.Debug()),
		)
	}

	return
}

func GetConflictOrError(
	el sku.ExternalLike,
) (f *fd.FD, err error) {
	var e *External
	ok := false

	if e, ok = el.(*External); !ok {
		err = errors.Errorf("expected store_fs.External but got %T", el)
		return
	}

	f = &e.fds.Conflict

	return
}

func GetObjectOrError(
	el sku.ExternalLike,
) (f *fd.FD, err error) {
	var e *External
	ok := false

	if e, ok = el.(*External); !ok {
		err = errors.Errorf("expected store_fs.External but got %T", el)
		return
	}

	f = &e.fds.Object

	return
}

func SetObjectOrError(
	el sku.ExternalLike,
	object *fd.FD,
) (err error) {
	var e *External
	ok := false

	if e, ok = el.(*External); !ok {
		err = errors.Errorf("expected store_fs.External but got %T", el)
		return
	}

	e.fds.Object.ResetWith(object)

	return
}

func GetBlobOrError(
	el sku.ExternalLike,
) (f *fd.FD, err error) {
	var e *External
	ok := false

	if e, ok = el.(*External); !ok {
		err = errors.Errorf("expected store_fs.External but got %T", el)
		return
	}

	f = &e.fds.Blob

	return
}

func SetBlobOrError(
	el sku.ExternalLike,
	blob *fd.FD,
) (err error) {
	var e *External
	ok := false

	if e, ok = el.(*External); !ok {
		err = errors.Errorf("expected store_fs.External but got %T", el)
		return
	}

	e.fds.Blob.ResetWith(blob)

	return
}

func UpdateDescriptionFromBlobs(
	el sku.ExternalLike,
) (err error) {
	var e *External
	ok := false

	if e, ok = el.(*External); !ok {
		err = errors.Errorf("expected store_fs.External but got %T", el)
		return
	}

	sorted := iter.ElementsSorted(
		e.fds.MutableSetLike,
		func(a, b *fd.FD) bool {
			return a.GetPath() < b.GetPath()
		},
	)

	for _, f := range sorted {
		desc := f.FileNameSansExt()

		if err = e.Transacted.Metadata.Description.Set(desc); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
