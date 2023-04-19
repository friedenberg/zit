package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type externalTextFormat[
	O Objekte[O],
	OPtr ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
] struct {
	standort         standort.Standort
	InlineChecker    kennung.InlineTypChecker
	AkteFactory      schnittstellen.AkteIOFactory
	AkteFormatter    script_config.RemoteScript
	TypError         error
	IncludeAkte      bool
	ExcludeMetadatei bool

	metadateiTextParser metadatei.TextParser
}

func MakeExternalTextFormatExcludeMetadatei[
	O Objekte[O],
	OPtr ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
](
	standort standort.Standort,
	inlineChecker kennung.InlineTypChecker,
	akteFactory schnittstellen.AkteIOFactory,
	akteFormatter script_config.RemoteScript,
) externalTextFormat[O, OPtr, K, KPtr] {
	return externalTextFormat[O, OPtr, K, KPtr]{
		standort:         standort,
		InlineChecker:    inlineChecker,
		AkteFactory:      akteFactory,
		AkteFormatter:    akteFormatter,
		IncludeAkte:      true,
		ExcludeMetadatei: true,
		metadateiTextParser: metadatei.MakeTextParser(
			akteFactory,
			akteFormatter,
		),
	}
}

func MakeExternalTextFormatIncludeAkte[
	O Objekte[O],
	OPtr ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
](
	standort standort.Standort,
	inlineChecker kennung.InlineTypChecker,
	akteFactory schnittstellen.AkteIOFactory,
	akteFormatter script_config.RemoteScript,
) externalTextFormat[O, OPtr, K, KPtr] {
	return externalTextFormat[O, OPtr, K, KPtr]{
		standort:      standort,
		InlineChecker: inlineChecker,
		AkteFactory:   akteFactory,
		AkteFormatter: akteFormatter,
		IncludeAkte:   true,
	}
}

func MakeExternalTextFormatAkteShaOnly[
	O Objekte[O],
	OPtr ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
](
	standort standort.Standort,
	akteFactory schnittstellen.AkteIOFactory,
	akteFormatter script_config.RemoteScript,
) externalTextFormat[O, OPtr, K, KPtr] {
	return externalTextFormat[O, OPtr, K, KPtr]{
		standort:      standort,
		AkteFactory:   akteFactory,
		AkteFormatter: akteFormatter,
	}
}

func (f externalTextFormat[O, OPtr, K, KPtr]) Format(
	w io.Writer,
	c *External[O, OPtr, K, KPtr],
) (n int64, err error) {
	inline := f.InlineChecker.IsInlineTyp(c.Objekte.GetMetadatei().GetTyp())

	var mtw io.WriterTo

	if !f.ExcludeMetadatei {
		mtw = format.MakeWriterToInterface[metadatei.TextFormatterContext](
			(metadatei.TextFormatter{
				IncludeAkteSha: !inline,
			}).Format,
			c,
		)
	}

	var wt io.WriterTo
	var ar sha.ReadCloser

	if inline {
		if ar, err = f.AkteFactory.AkteReader(
			c.Objekte.GetMetadatei().AkteSha,
		); err != nil {
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

	mw := metadatei.Writer{
		Metadatei: mtw,
		Akte:      wt,
	}

	if n, err = mw.WriteTo(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
