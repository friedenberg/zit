package zettel_external

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/metadatei_io"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type fileEncoder struct {
	mode int
	perm os.FileMode
	arf  schnittstellen.AkteReaderFactory
	ic   typ.InlineChecker
}

func MakeFileEncoder(
	arf schnittstellen.AkteReaderFactory,
	ic typ.InlineChecker,
) fileEncoder {
	return fileEncoder{
		mode: os.O_WRONLY | os.O_CREATE | os.O_EXCL | os.O_APPEND,
		perm: 0o666,
		arf:  arf,
		ic:   ic,
	}
}

func MakeFileEncoderJustOpen(
	arf schnittstellen.AkteReaderFactory,
	ic typ.InlineChecker,
) fileEncoder {
	return fileEncoder{
		mode: os.O_WRONLY | os.O_EXCL | os.O_APPEND,
		perm: 0o666,
		arf:  arf,
		ic:   ic,
	}
}

func (e *fileEncoder) openOrCreate(p string) (*os.File, error) {
	return files.OpenFile(p, e.mode, e.perm)
}

func (e *fileEncoder) EncodeObjekte(
	z *zettel.Objekte,
	objektePath string,
	aktePath string,
) (err error) {
	inline := e.ic.IsInlineTyp(z.Typ)

	mtw := zettel.TextMetadateiFormatter{
		IncludeAkteSha: !inline,
	}

	var ar sha.ReadCloser

	if ar, err = e.arf.AkteReader(z.Akte); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	mw := metadatei_io.Writer{
		Metadatei: format.MakeWriterTo2(
			mtw.Format,
			&zettel.Metadatei{
				Objekte:  *z,
				AktePath: aktePath,
			},
		),
	}

	switch {
	case aktePath != "" && objektePath != "":
		var fAkte, fZettel *os.File

		if fAkte, err = e.openOrCreate(
			aktePath,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, fAkte.Close)

		if fZettel, err = e.openOrCreate(
			objektePath,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, fZettel.Close)

		if _, err = mw.WriteTo(fZettel); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = io.Copy(fAkte, ar); err != nil {
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

		defer errors.Deferred(&err, fAkte.Close)

		if _, err = io.Copy(fAkte, ar); err != nil {
			err = errors.Wrap(err)
			return
		}

	case objektePath != "":
		if inline {
			mw.Akte = ar
		}

		var fZettel *os.File

		if fZettel, err = e.openOrCreate(
			objektePath,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, fZettel.Close)

		if _, err = mw.WriteTo(fZettel); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (e *fileEncoder) Encode(z *zettel.External) (err error) {
	return e.EncodeObjekte(
		&z.Objekte,
		z.GetObjekteFD().Path,
		z.GetAkteFD().Path,
	)
}
