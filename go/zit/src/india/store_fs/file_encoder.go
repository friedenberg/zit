package store_fs

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
)

type FileEncoder interface {
	Encode(options checkout_options.TextFormatterOptions, z *External) (err error)
}

type fileEncoder struct {
	mode int
	perm os.FileMode
	arf  interfaces.BlobIOFactory
	ic   ids.InlineTypeChecker
}

func MakeFileEncoder(
	arf interfaces.BlobIOFactory,
	ic ids.InlineTypeChecker,
) *fileEncoder {
	return &fileEncoder{
		mode: os.O_WRONLY | os.O_CREATE | os.O_TRUNC,
		perm: 0o666,
		arf:  arf,
		ic:   ic,
	}
}

func MakeFileEncoderJustOpen(
	arf interfaces.BlobIOFactory,
	ic ids.InlineTypeChecker,
) fileEncoder {
	return fileEncoder{
		mode: os.O_WRONLY | os.O_TRUNC,
		perm: 0o666,
		arf:  arf,
		ic:   ic,
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

func (e *fileEncoder) EncodeObjekte(
	options checkout_options.TextFormatterOptions,
	z *External,
	objektePath string,
	aktePath string,
) (err error) {
	inline := e.ic.IsInlineType(z.GetType())

	var ar sha.ReadCloser

	if ar, err = e.arf.BlobReader(z.GetBlobSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	switch {
	case aktePath != "" && objektePath != "":
		mtw := object_metadata.MakeTextFormatterMetadateiAktePath(
			options,
			e.arf,
			nil,
		)

		var fAkte, fZettel *os.File

		{
			if fAkte, err = e.openOrCreate(
				aktePath,
			); err != nil {
				if errors.IsExist(err) {
					var aw sha.WriteCloser

					if aw, err = e.arf.BlobWriter(); err != nil {
						err = errors.Wrap(err)
						return
					}

					defer errors.DeferredCloser(&err, aw)

					if _, err = io.Copy(aw, fAkte); err != nil {
						err = errors.Wrap(err)
						return
					}

				} else {
					err = errors.Wrap(err)
					return
				}
			}

			defer errors.DeferredCloser(&err, fAkte)

			if _, err = io.Copy(fAkte, ar); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if fZettel, err = e.openOrCreate(
			objektePath,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, fZettel)

		if _, err = mtw.FormatMetadata(fZettel, z); err != nil {
			err = errors.Wrap(err)
			return
		}

	case aktePath != "":
		var fAkte *os.File

		if fAkte, err = e.openOrCreate(
			aktePath,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, fAkte)

		if _, err = io.Copy(fAkte, ar); err != nil {
			err = errors.Wrap(err)
			return
		}

	case objektePath != "":
		var mtw object_metadata.TextFormatter

		if inline {
			mtw = object_metadata.MakeTextFormatterMetadataInlineBlob(
				options,
				e.arf,
				nil,
			)
		} else {
			mtw = object_metadata.MakeTextFormatterMetadataOnly(
				options,
				e.arf,
				nil,
			)
		}

		var fZettel *os.File

		if fZettel, err = e.openOrCreate(
			objektePath,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, fZettel)

		if _, err = mtw.FormatMetadata(fZettel, z); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (e *fileEncoder) Encode(
	options checkout_options.TextFormatterOptions,
	z *External,
) (err error) {
	return e.EncodeObjekte(
		options,
		z,
		z.GetObjekteFD().GetPath(),
		z.GetAkteFD().GetPath(),
	)
}
