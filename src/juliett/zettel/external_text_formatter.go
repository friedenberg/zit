package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/metadatei_io"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/typ"
)

type externalTextFormatter struct {
	standort         standort.Standort
	InlineChecker    typ.InlineChecker
	AkteFactory      schnittstellen.AkteIOFactory
	AkteFormatter    erworben.RemoteScript
	TypError         error
	IncludeAkte      bool
	ExcludeMetadatei bool

	objekteTextParser
}

func MakeExternalTextFormatterExcludeMetadatei(
	standort standort.Standort,
	inlineChecker typ.InlineChecker,
	akteFactory schnittstellen.AkteIOFactory,
	akteFormatter erworben.RemoteScript,
) externalTextFormatter {
	return externalTextFormatter{
		standort:         standort,
		InlineChecker:    inlineChecker,
		AkteFactory:      akteFactory,
		AkteFormatter:    akteFormatter,
		IncludeAkte:      true,
		ExcludeMetadatei: true,
		objekteTextParser: MakeObjekteTextParser(
			akteFactory,
			akteFormatter,
		),
	}
}

func MakeExternalTextFormatterIncludeAkte(
	standort standort.Standort,
	inlineChecker typ.InlineChecker,
	akteFactory schnittstellen.AkteIOFactory,
	akteFormatter erworben.RemoteScript,
) externalTextFormatter {
	return externalTextFormatter{
		standort:      standort,
		InlineChecker: inlineChecker,
		AkteFactory:   akteFactory,
		AkteFormatter: akteFormatter,
		IncludeAkte:   true,
	}
}

func MakeExternalTextFormatterAkteShaOnly(
	standort standort.Standort,
	akteFactory schnittstellen.AkteIOFactory,
	akteFormatter erworben.RemoteScript,
) externalTextFormatter {
	return externalTextFormatter{
		standort:      standort,
		AkteFactory:   akteFactory,
		AkteFormatter: akteFormatter,
	}
}

func (f externalTextFormatter) Format(
	w io.Writer,
	c *External,
) (n int64, err error) {
	inline := f.InlineChecker.IsInlineTyp(c.Objekte.Typ)

	var mtw io.WriterTo

	if !f.ExcludeMetadatei {
		mtw = format.MakeWriterTo2(
			(&TextMetadateiFormatter{
				IncludeAkteSha: !inline,
			}).Format,
			&Metadatei{
				Objekte: c.Objekte,
			},
		)
	}

	var wt io.WriterTo
	var ar sha.ReadCloser

	if inline {
		if ar, err = f.AkteFactory.AkteReader(c.Objekte.Akte); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, ar.Close)

		wt = ar
	}

	if f.AkteFormatter != nil {
		if wt, err = script_config.MakeWriterToWithStdin(
			f.AkteFormatter,
			map[string]string{
				"ZIT_BIN": f.standort.Executable(),
			},
			ar,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	mw := metadatei_io.Writer{
		Metadatei: mtw,
		Akte:      wt,
	}

	if n, err = mw.WriteTo(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
