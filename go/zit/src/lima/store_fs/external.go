package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) GetCheckoutModeOrError(
	el sku.ExternalLike,
	originalMode checkout_mode.Mode,
) (m checkout_mode.Mode, err error) {
	var fds *Item

	if fds, err = s.ReadFSItemFromExternal(el); err != nil {
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

	if fds, err = s.ReadFSItemFromExternal(el); err != nil {
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

	if fds, err = s.ReadFSItemFromExternal(el); err != nil {
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

	if fds, err = s.ReadFSItemFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	fds.Object.ResetWith(object)

	if err = s.WriteFSItemToExternal(fds, el); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) GetBlobOrError(
	el sku.ExternalLike,
) (f *fd.FD, err error) {
	var fds *Item

	if fds, err = s.ReadFSItemFromExternal(el); err != nil {
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

	if fds, err = s.ReadFSItemFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	fds.Blob.ResetWith(blob)

	if err = s.WriteFSItemToExternal(fds, el); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) UpdateDescriptionFromBlobs(
	el sku.ExternalLike,
) (err error) {
	var fds *Item

	if fds, err = s.ReadFSItemFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	sorted := quiter.ElementsSorted(
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

	if err = s.WriteFSItemToExternal(fds, el); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
