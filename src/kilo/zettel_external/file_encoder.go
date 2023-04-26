package zettel_external

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type fileEncoder struct {
	mode int
	perm os.FileMode
	arf  schnittstellen.AkteIOFactory
	ic   kennung.InlineTypChecker
}

func MakeFileEncoder(
	arf schnittstellen.AkteIOFactory,
	ic kennung.InlineTypChecker,
) fileEncoder {
	return fileEncoder{
		mode: os.O_WRONLY | os.O_CREATE | os.O_EXCL | os.O_APPEND,
		perm: 0o666,
		arf:  arf,
		ic:   ic,
	}
}

func MakeFileEncoderJustOpen(
	arf schnittstellen.AkteIOFactory,
	ic kennung.InlineTypChecker,
) fileEncoder {
	return fileEncoder{
		mode: os.O_WRONLY | os.O_EXCL | os.O_APPEND,
		perm: 0o666,
		arf:  arf,
		ic:   ic,
	}
}

func (e *fileEncoder) openOrCreate(p string) (f *os.File, err error) {
	if f, err = files.OpenFile(p, e.mode, e.perm); err != nil {
		if errors.IsExist(err) {
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
	z *zettel.External,
	objektePath string,
	aktePath string,
) (err error) {
	inline := e.ic.IsInlineTyp(z.GetTyp())

	var ar sha.ReadCloser

	if ar, err = e.arf.AkteReader(z.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	switch {
	case aktePath != "" && objektePath != "":
		mtw := metadatei.MakeTextFormatterMetadateiAktePath(
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

					if aw, err = e.arf.AkteWriter(); err != nil {
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
			} else {
				defer errors.DeferredCloser(&err, fAkte)

				if _, err = io.Copy(fAkte, ar); err != nil {
					err = errors.Wrap(err)
					return
				}
			}
		}

		if fZettel, err = e.openOrCreate(
			objektePath,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, fZettel)

		if _, err = mtw.FormatMetadatei(fZettel, z); err != nil {
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
		var mtw metadatei.TextFormatter

		if inline {
			mtw = metadatei.MakeTextFormatterMetadateiInlineAkte(
				e.arf,
				nil,
			)
		} else {
			mtw = metadatei.MakeTextFormatterMetadateiOnly(
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

		if _, err = mtw.FormatMetadatei(fZettel, z); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (e *fileEncoder) Encode(z *zettel.External) (err error) {
	return e.EncodeObjekte(
		z,
		z.GetObjekteFD().Path,
		z.GetAkteFD().Path,
	)
}
