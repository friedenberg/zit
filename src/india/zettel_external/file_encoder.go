package zettel_external

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/format"
	"github.com/friedenberg/zit/src/foxtrot/metadatei_io"
	"github.com/friedenberg/zit/src/foxtrot/sha"
	"github.com/friedenberg/zit/src/india/typ"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

type fileEncoder struct {
	arf gattung.AkteReaderFactory
	ic  typ.InlineChecker
}

func MakeFileEncoder(
	arf gattung.AkteReaderFactory,
	ic typ.InlineChecker,
) fileEncoder {
	return fileEncoder{
		arf: arf,
		ic:  ic,
	}
}

func (e *fileEncoder) Encode(z *Zettel) (err error) {
	inline := e.ic.IsInlineTyp(z.Objekte.Typ)

	mtw := zettel.TextMetadateiFormatter{
		IncludeAkteSha: !inline,
	}

	var ar sha.ReadCloser

	if ar, err = e.arf.AkteReader(z.Objekte.Akte); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, ar.Close)

	mw := metadatei_io.Writer{
		Metadatei: format.MakeWriterTo2(
			mtw.Format,
			&zettel.Metadatei{
				Objekte:  z.Objekte,
				AktePath: z.AkteFD.Path,
			},
		),
	}

	switch {
	case z.AkteFD.Path != "" && z.ZettelFD.Path != "":
		var fAkte, fZettel *os.File

		if fAkte, err = files.CreateExclusiveWriteOnly(z.AkteFD.Path); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, fAkte.Close)

		if fZettel, err = files.CreateExclusiveWriteOnly(z.ZettelFD.Path); err != nil {
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

	case z.AkteFD.Path != "":
		var fAkte *os.File

		if fAkte, err = files.CreateExclusiveWriteOnly(z.AkteFD.Path); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, fAkte.Close)

		if _, err = io.Copy(fAkte, ar); err != nil {
			err = errors.Wrap(err)
			return
		}

	case z.ZettelFD.Path != "":
		if inline {
			mw.Akte = ar
		}

		var fZettel *os.File

		if fZettel, err = files.CreateExclusiveWriteOnly(z.ZettelFD.Path); err != nil {
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
