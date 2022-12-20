package zettel

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/metadatei_io"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/india/typ"
)

const MetadateiBoundary = metadatei_io.Boundary

type objekteTextFormatter struct {
	InlineChecker typ.InlineChecker
	AkteFactory   gattung.AkteIOFactory
	AkteFormatter konfig.RemoteScript
	TypError      error
	IncludeAkte   bool
}

func MakeObjekteTextFormatterIncludeAkte(
	inlineChecker typ.InlineChecker,
	akteFactory gattung.AkteIOFactory,
	akteFormatter konfig.RemoteScript,
) objekteTextFormatter {
	return objekteTextFormatter{
		InlineChecker: inlineChecker,
		AkteFactory:   akteFactory,
		AkteFormatter: akteFormatter,
		IncludeAkte:   true,
	}
}

func MakeObjekteTextFormatterAkteShaOnly(
	akteFactory gattung.AkteIOFactory,
	akteFormatter konfig.RemoteScript,
) objekteTextFormatter {
	return objekteTextFormatter{
		AkteFactory:   akteFactory,
		AkteFormatter: akteFormatter,
	}
}

// TODO switch to three different formats
// metadatei, zettel-akte-external, zettel-akte-inline
func (f objekteTextFormatter) Format(
	w io.Writer,
	c *FormatContextWrite,
) (n int64, err error) {
	inline := f.InlineChecker.IsInlineTyp(c.Zettel.Typ)

	mtw := TextMetadateiFormatter{
		IncludeAkteSha: !inline,
	}

	var ar sha.ReadCloser

	if inline {
		if ar, err = f.AkteFactory.AkteReader(c.Zettel.Akte); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, ar.Close)
	}

	mw := metadatei_io.Writer{
		Metadatei: format.MakeWriterTo2(
			mtw.Format,
			&Metadatei{
				Objekte: c.Zettel,
			},
		),
		Akte: ar,
	}

	if n, err = mw.WriteTo(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f objekteTextFormatter) writeToExternalAkte(
	w1 io.Writer,
	c *FormatContextWrite,
) (n int64, err error) {
	w := format.NewWriter()

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
