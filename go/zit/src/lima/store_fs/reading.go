package store_fs

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) UpdateTransacted(z *sku.Transacted) (err error) {
	e, ok := s.Get(&z.ObjectId)

	if !ok {
		return
	}

	var e2 *sku.External

	if e2, err = s.ReadExternalFromItem(
		sku.CommitOptions{
			Mode: objekte_mode.ModeUpdateTai,
		},
		e,
		z,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sku.Resetter.ResetWith(z, e2)

	return
}

func (s *Store) readOneExternalInto(
	o *sku.CommitOptions,
	i *Item,
	t *sku.Transacted,
	e *sku.External,
) (err error) {
	if err = s.WriteFSItemToExternal(i, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	if t != nil {
		e.ObjectId.ResetWith(&t.ObjectId)
	}

	var m checkout_mode.Mode

	if m, err = i.GetCheckoutModeOrError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var t1 *sku.Transacted

	if t != nil {
		t1 = t
	}

	switch m {
	case checkout_mode.BlobOnly:
		if err = s.ReadOneExternalBlob(e, t1, i); err != nil {
			err = errors.Wrap(err)
			return
		}

	case checkout_mode.MetadataOnly, checkout_mode.MetadataAndBlob:
		if i.Object.IsStdin() {
			if err = s.ReadOneExternalObjectReader(os.Stdin, e); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else {
			if err = s.readOneExternalObject(e, t1, i); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	case checkout_mode.BlobRecognized:
		object_metadata.Resetter.ResetWith(
			e.GetMetadata(),
			t1.GetMetadata(),
		)

	default:
		panic(checkout_mode.MakeErrInvalidCheckoutModeMode(m))
	}

	if !i.Blob.IsEmpty() {
		blobFD := &i.Blob
		ext := blobFD.ExtSansDot()
		typFromExtension := s.config.GetTypeStringFromExtension(ext)

		if typFromExtension == "" {
			typFromExtension = ext
		}

		if typFromExtension != "" {
			if err = e.Metadata.Type.Set(typFromExtension); err != nil {
				err = errors.Wrapf(err, "Path: %s", blobFD.GetPath())
				return
			}
		}
	}

	if o.Clock == nil {
		o.Clock = i
	}

	if err = s.WriteFSItemToExternal(i, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) readOneExternalObject(
	e *sku.External,
	t *sku.Transacted,
	i *Item,
) (err error) {
	if t != nil {
		object_metadata.Resetter.ResetWith(
			e.GetMetadata(),
			t.GetMetadata(),
		)
	}

	var f *os.File

	if f, err = files.Open(i.Object.GetPath()); err != nil {
		err = errors.Wrapf(err, "Item: %s", i.Debug())
		return
	}

	defer errors.DeferredCloser(&err, f)

	if err = s.ReadOneExternalObjectReader(f, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadOneExternalObjectReader(
	r io.Reader,
	e *sku.External,
) (err error) {
	if _, err = s.metadataTextParser.ParseMetadata(r, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadOneExternalBlob(
	e *sku.External,
	t *sku.Transacted,
	i *Item,
) (err error) {
	object_metadata.Resetter.ResetWith(&e.Metadata, t.GetMetadata())

	// TODO use cache
	{
		var aw sha.WriteCloser

		if aw, err = s.fs_home.BlobWriter(); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, aw)

		var f *os.File

		if f, err = files.OpenExclusiveReadOnly(
			i.Blob.GetPath(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, f)

		if _, err = io.Copy(aw, f); err != nil {
			err = errors.Wrap(err)
			return
		}

		e.GetMetadata().Blob.SetShaLike(aw)
	}

	return
}
