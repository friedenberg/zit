package store_fs

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) HydrateExternalFromItem(
	o sku.CommitOptions,
	item *sku.FSItem,
	internal *sku.Transacted,
	external *sku.Transacted,
) (err error) {
	o.Del(object_mode.ModeApplyProto)

	if err = s.readOneExternalInto(&o, item, internal, external); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.externalStoreSupplies.FuncCommit(
		external,
		o,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) readOneExternalInto(
	o *sku.CommitOptions,
	item *sku.FSItem,
	internal *sku.Transacted,
	external *sku.Transacted,
) (err error) {
	if err = s.WriteFSItemToExternal(item, external); err != nil {
		err = errors.Wrap(err)
		return
	}

	if internal != nil {
		external.ObjectId.ResetWith(&internal.ObjectId)
	}

	var m checkout_mode.Mode

	if m, err = item.GetCheckoutModeOrError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	switch m {
	case checkout_mode.BlobOnly:
		if err = s.readOneExternalBlob(external, internal, item); err != nil {
			err = errors.Wrap(err)
			return
		}

	case checkout_mode.MetadataOnly, checkout_mode.MetadataAndBlob:
		if item.Object.IsStdin() {
			if err = s.ReadOneExternalObjectReader(os.Stdin, external); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else {
			if err = s.readOneExternalObject(external, internal, item); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	case checkout_mode.BlobRecognized:
		object_metadata.Resetter.ResetWith(
			external.GetMetadata(),
			internal.GetMetadata(),
		)

	default:
		panic(checkout_mode.MakeErrInvalidCheckoutModeMode(m))
	}

	// if !i.Blob.IsEmpty() {
	// 	blobFD := &i.Blob
	// 	ext := blobFD.ExtSansDot()
	// 	typFromExtension := s.config.GetTypeStringFromExtension(ext)

	// 	if typFromExtension == "" {
	// 		typFromExtension = ext
	// 	}

	// 	if typFromExtension != "" {
	// 		if err = e.Metadata.Type.Set(typFromExtension); err != nil {
	// 			err = errors.Wrapf(err, "Path: %s", blobFD.GetPath())
	// 			return
	// 		}
	// 	}
	// }

	if o.Clock == nil {
		o.Clock = item
	}

	if err = s.WriteFSItemToExternal(item, external); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) readOneExternalObject(
	external *sku.Transacted,
	internal *sku.Transacted,
	item *sku.FSItem,
) (err error) {
	if internal != nil {
		object_metadata.Resetter.ResetWith(
			external.GetMetadata(),
			internal.GetMetadata(),
		)
	}

	var f *os.File

	if f, err = files.Open(item.Object.GetPath()); err != nil {
		err = errors.Wrapf(err, "Item: %s", item.Debug())
		return
	}

	defer errors.DeferredCloser(&err, f)

	if err = s.ReadOneExternalObjectReader(f, external); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) readOneExternalBlob(
	external *sku.Transacted,
	internal *sku.Transacted,
	item *sku.FSItem,
) (err error) {
	object_metadata.Resetter.ResetWith(&external.Metadata, internal.GetMetadata())

	// TODO use cache
	{
		var aw sha.WriteCloser

		if aw, err = s.dirLayout.BlobWriter(); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, aw)

		var f *os.File

		if f, err = files.OpenExclusiveReadOnly(
			item.Blob.GetPath(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, f)

		if _, err = io.Copy(aw, f); err != nil {
			err = errors.Wrap(err)
			return
		}

		external.GetMetadata().Blob.SetShaLike(aw)
	}

	return
}
