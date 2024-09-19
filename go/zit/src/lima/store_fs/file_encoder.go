package store_fs

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type FileEncoder interface {
	Encode(
		options checkout_options.TextFormatterOptions,
		z *sku.External,
		i *Item,
	) (err error)
}

type fileEncoder struct {
	mode    int
	perm    os.FileMode
	fs_home fs_home.Home
	ic      ids.InlineTypeChecker

	object_metadata.TextFormatterFamily
}

func MakeFileEncoder(
	fs_home fs_home.Home,
	ic ids.InlineTypeChecker,
) *fileEncoder {
	return &fileEncoder{
		mode:    os.O_WRONLY | os.O_CREATE | os.O_TRUNC,
		perm:    0o666,
		fs_home: fs_home,
		ic:      ic,
		TextFormatterFamily: object_metadata.MakeTextFormatterFamily(
			fs_home,
			nil,
		),
	}
}

func (e *fileEncoder) openOrCreate(p string) (f *os.File, err error) {
	if f, err = files.OpenFile(p, e.mode, e.perm); err != nil {
		err = errors.Wrap(err)

		if errors.IsExist(err) {
			// err = nil
			var err2 error

			if f, err2 = files.OpenExclusiveReadOnly(p); err2 != nil {
				err = errors.Wrap(err2)
			}
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	return
}

func (e *fileEncoder) EncodeObject(
	options checkout_options.TextFormatterOptions,
	z *sku.External,
	objectPath string,
	blobPath string,
) (err error) {
	ctx := object_metadata.TextFormatterContext{
		PersistentFormatterContext: z.GetSku(),
		TextFormatterOptions:       options,
	}

	inline := e.ic.IsInlineType(z.GetType())

	var ar sha.ReadCloser

	if ar, err = e.fs_home.BlobReader(z.Transacted.GetBlobSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	switch {
	case blobPath != "" && objectPath != "":
		var fBlob, fZettel *os.File

		{
			if fBlob, err = e.openOrCreate(
				blobPath,
			); err != nil {
				if errors.IsExist(err) {
					var aw sha.WriteCloser

					if aw, err = e.fs_home.BlobWriter(); err != nil {
						err = errors.Wrap(err)
						return
					}

					defer errors.DeferredCloser(&err, aw)

					if _, err = io.Copy(aw, fBlob); err != nil {
						err = errors.Wrap(err)
						return
					}

				} else {
					err = errors.Wrap(err)
					return
				}
			}

			defer errors.DeferredCloser(&err, fBlob)

			if _, err = io.Copy(fBlob, ar); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if fZettel, err = e.openOrCreate(
			objectPath,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, fZettel)

		if _, err = e.BlobPath.FormatMetadata(fZettel, ctx); err != nil {
			err = errors.Wrap(err)
			return
		}

	case blobPath != "":
		var fBlob *os.File

		if fBlob, err = e.openOrCreate(
			blobPath,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, fBlob)

		if _, err = io.Copy(fBlob, ar); err != nil {
			err = errors.Wrap(err)
			return
		}

	case objectPath != "":
		var mtw object_metadata.TextFormatter

		if inline {
			mtw = e.InlineBlob
		} else {
			mtw = e.MetadataOnly
		}

		var fZettel *os.File

		if fZettel, err = e.openOrCreate(
			objectPath,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, fZettel)

		if _, err = mtw.FormatMetadata(fZettel, ctx); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (e *fileEncoder) Encode(
	options checkout_options.TextFormatterOptions,
	z *sku.External,
	i *Item,
) (err error) {
	return e.EncodeObject(
		options,
		z,
		i.Object.GetPath(),
		i.Blob.GetPath(),
	)
}
