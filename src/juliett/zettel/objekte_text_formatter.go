package zettel

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/metadatei_io"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/echo/standort"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/typ"
)

const MetadateiBoundary = metadatei_io.Boundary

type objekteTextFormatter struct {
	standort         standort.Standort
	InlineChecker    typ.InlineChecker
	AkteFactory      schnittstellen.AkteIOFactory
	AkteFormatter    erworben.RemoteScript
	TypError         error
	IncludeAkte      bool
	ExcludeMetadatei bool
}

func MakeObjekteTextFormatterExcludeMetadatei(
	standort standort.Standort,
	inlineChecker typ.InlineChecker,
	akteFactory schnittstellen.AkteIOFactory,
	akteFormatter erworben.RemoteScript,
) objekteTextFormatter {
	return objekteTextFormatter{
		standort:         standort,
		InlineChecker:    inlineChecker,
		AkteFactory:      akteFactory,
		AkteFormatter:    akteFormatter,
		IncludeAkte:      true,
		ExcludeMetadatei: true,
	}
}

func MakeObjekteTextFormatterIncludeAkte(
	standort standort.Standort,
	inlineChecker typ.InlineChecker,
	akteFactory schnittstellen.AkteIOFactory,
	akteFormatter erworben.RemoteScript,
) objekteTextFormatter {
	return objekteTextFormatter{
		standort:      standort,
		InlineChecker: inlineChecker,
		AkteFactory:   akteFactory,
		AkteFormatter: akteFormatter,
		IncludeAkte:   true,
	}
}

func MakeObjekteTextFormatterAkteShaOnly(
	standort standort.Standort,
	akteFactory schnittstellen.AkteIOFactory,
	akteFormatter erworben.RemoteScript,
) objekteTextFormatter {
	return objekteTextFormatter{
		standort:      standort,
		AkteFactory:   akteFactory,
		AkteFormatter: akteFormatter,
	}
}

func (f objekteTextFormatter) Format(
	w io.Writer,
	c *ObjekteFormatterContext,
) (n int64, err error) {
	inline := f.InlineChecker.IsInlineTyp(c.Zettel.Typ)

	var mtw io.WriterTo

	if !f.ExcludeMetadatei {
		mtw = format.MakeWriterTo2(
			(&TextMetadateiFormatter{
				IncludeAkteSha: !inline,
			}).Format,
			&Metadatei{
				Objekte: c.Zettel,
			},
		)
	}

	var wt io.WriterTo
	var ar sha.ReadCloser

	if inline {
		if ar, err = f.AkteFactory.AkteReader(c.Zettel.Akte); err != nil {
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

func (f objekteTextFormatter) writeToExternalAkte(
	w1 io.Writer,
	c *ObjekteFormatterContext) (n int64, err error) {
	w := format.NewLineWriter()

	w.WriteLines(
		MetadateiBoundary,
		fmt.Sprintf("# %s", c.Zettel.Bezeichnung),
	)

	for _, e := range c.Zettel.Etiketten.Sorted() {
		w.WriteFormat("- %s", e)
	}

	if strings.Index(c.ExternalAktePath, "\n") != -1 {
		panic(errors.Errorf("ExternalAktePath contains newline: %q", c.ExternalAktePath))
	}

	w.WriteLines(
		fmt.Sprintf("! %s", c.ExternalAktePath),
	)

	w.WriteLines(
		MetadateiBoundary,
	)

	n, err = w.WriteTo(w1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	var ar io.ReadCloser

	if f.AkteFactory == nil {
		err = errors.Errorf("akte reader factory is nil")
		return
	}

	if ar, err = f.AkteFactory.AkteReader(c.Zettel.Akte); err != nil {
		err = errors.Wrap(err)
		return
	}

	if ar == nil {
		err = errors.Errorf("akte reader is nil")
		return
	}

	defer errors.Deferred(&err, ar.Close)

	var file *os.File

	if file, err = files.Create(c.ExternalAktePath); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, file.Close)

	var n1 int64

	n1, err = io.Copy(file, ar)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
