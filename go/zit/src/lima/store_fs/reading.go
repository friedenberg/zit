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

func (s *Store) ReadOneExternal(
	o *sku.CommitOptions,
	em *FDSet,
	t *sku.Transacted,
) (e *External, err error) {
	e = GetExternalPool().Get()

	if err = s.ReadOneExternalInto(o, em, t, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) UpdateTransacted(z *sku.Transacted) (err error) {
	e, ok := s.Get(&z.ObjectId)

	if !ok {
		return
	}

	var e2 *External

	if e2, err = s.ReadExternalFromObjectIdFDPair(
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

func (s *Store) ReadOneExternalInto(
	o *sku.CommitOptions,
	em *FDSet,
	t *sku.Transacted,
	e *External,
) (err error) {
	if err = e.ResetWithExternalMaybe(em); err != nil {
		err = errors.Wrap(err)
		return
	}

	if t != nil {
		e.Transacted.ObjectId.ResetWith(&t.ObjectId)
	}

	var m checkout_mode.Mode

	if m, err = em.GetCheckoutModeOrError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var t1 *sku.Transacted

	if t != nil {
		t1 = t
	}

	switch m {
	case checkout_mode.BlobOnly:
		if err = s.ReadOneExternalBlob(e, t1); err != nil {
			err = errors.Wrap(err)
			return
		}

	case checkout_mode.MetadataOnly, checkout_mode.MetadataAndBlob:
		if e.fds.Object.IsStdin() {
			if err = s.ReadOneExternalObjectReader(os.Stdin, e); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else {
			if err = s.ReadOneExternalObject(e, t1); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	case checkout_mode.BlobRecognized:
		object_metadata.Resetter.ResetWith(
      e.Transacted.GetMetadata(),
      t1.GetMetadata(),
    )

	default:
		panic(checkout_mode.MakeErrInvalidCheckoutModeMode(m))
	}

	if !e.fds.Blob.IsEmpty() {
		blobFD := &e.fds.Blob
		ext := blobFD.ExtSansDot()
		typFromExtension := s.config.GetTypeStringFromExtension(ext)

		if typFromExtension == "" {
			typFromExtension = ext
		}

		if typFromExtension != "" {
			if err = e.Transacted.Metadata.Type.Set(typFromExtension); err != nil {
				err = errors.Wrapf(err, "Path: %s", blobFD.GetPath())
				return
			}
		}
	}

	if o.Clock == nil {
		o.Clock = &e.fds
	}

	return
}

func (s *Store) ReadOneExternalObject(
	e *External,
	t *sku.Transacted,
) (err error) {
	if t != nil {
		object_metadata.Resetter.ResetWith(
      e.Transacted.GetMetadata(),
      t.GetMetadata(),
    )
	}

	var f *os.File

	if f, err = files.Open(e.GetObjectFD().GetPath()); err != nil {
		err = errors.Wrap(err)
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
	e *External,
) (err error) {
	if _, err = s.metadataTextParser.ParseMetadata(r, &e.Transacted); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadOneExternalBlob(
	e *External,
	t *sku.Transacted,
) (err error) {
	object_metadata.Resetter.ResetWith(&e.Transacted.Metadata, t.GetMetadata())

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
			e.GetBlobFD().GetPath(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, f)

		if _, err = io.Copy(aw, f); err != nil {
			err = errors.Wrap(err)
			return
		}

		e.Transacted.GetMetadata().Blob.SetShaLike(aw)
	}

	return
}
