package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) GetCheckoutMode(
	el sku.ExternalLike,
) (m checkout_mode.Mode, err error) {
	var fds *sku.FSItem

	if fds, err = s.ReadFSItemFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	m = fds.GetCheckoutMode()

	return
}

func (s *Store) GetCheckoutModeOrError(
	el sku.ExternalLike,
) (m checkout_mode.Mode, err error) {
	var fds *sku.FSItem

	if fds, err = s.ReadFSItemFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	if m, err = fds.GetCheckoutModeOrError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) GetConflictOrError(
	el sku.ExternalLike,
) (f *fd.FD, err error) {
	var fds *sku.FSItem

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
	var fds *sku.FSItem

	if fds, err = s.ReadFSItemFromExternal(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	f = &fds.Object

	return
}

func (s *Store) UpdateTransactedFromBlobs(
	el sku.ExternalLike,
) (err error) {
	var fds *sku.FSItem

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

	if !fds.Blob.IsEmpty() {
		blobFD := &fds.Blob
		ext := blobFD.ExtSansDot()
		typFromExtension := s.config.GetTypeStringFromExtension(ext)

		if typFromExtension == "" {
			typFromExtension = ext
		}

		if typFromExtension != "" {
			if err = el.GetSku().Metadata.Type.Set(typFromExtension); err != nil {
				err = errors.Wrapf(err, "Path: %s", blobFD.GetPath())
				return
			}
		}
	}

	if err = s.WriteFSItemToExternal(fds, el); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
