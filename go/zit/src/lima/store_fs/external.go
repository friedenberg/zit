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
}

func (t *External) GetSkuExternalLike() sku.ExternalLike {
	return t
}

func (a *External) Clone() sku.ExternalLike {
	b := GetExternalPool().Get()
	sku.ExternalResetter.ResetWith(&b.External, &a.External)
	return b
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

	return
}

func (o *External) GetKey() string {
	return fmt.Sprintf("%s.%s", o.GetGenre(), o.GetObjectId())
}

func (s *Store) GetCheckoutModeOrError(
	el sku.ExternalLike,
	originalMode checkout_mode.Mode,
) (m checkout_mode.Mode, err error) {
	var fds *Item

	if fds, err = s.ReadFromExternal(el); err != nil {
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
			errors.Errorf(
				"all FD's are empty: %s. Fields: %#v",
				fds.Debug(),
				el.GetSku().Metadata.Fields,
			),
		)
	}

	return
}

func (s *Store) GetConflictOrError(
	el sku.ExternalLike,
) (f *fd.FD, err error) {
	var fds *Item

	if fds, err = s.ReadFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	f = &fds.Conflict

	return
}

func (s *Store) GetObjectOrError(
	el sku.ExternalLike,
) (f *fd.FD, err error) {
	var fds *Item

	if fds, err = s.ReadFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	f = &fds.Object

	return
}

func (s *Store) SetObjectOrError(
	el sku.ExternalLike,
	object *fd.FD,
) (err error) {
	var fds *Item

	if fds, err = s.ReadFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	fds.Object.ResetWith(object)

	if err = s.WriteToExternal(fds, el); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) GetBlobOrError(
	el sku.ExternalLike,
) (f *fd.FD, err error) {
	var fds *Item

	if fds, err = s.ReadFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	f = &fds.Blob

	return
}

func (s *Store) SetBlobOrError(
	el sku.ExternalLike,
	blob *fd.FD,
) (err error) {
	var fds *Item

	if fds, err = s.ReadFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	fds.Blob.ResetWith(blob)

	if err = s.WriteToExternal(fds, el); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) UpdateDescriptionFromBlobs(
	el sku.ExternalLike,
) (err error) {
	var fds *Item

	if fds, err = s.ReadFromExternal(el); err != nil {
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

	if err = s.WriteToExternal(fds, el); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
