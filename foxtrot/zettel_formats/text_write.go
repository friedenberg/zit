package zettel_formats

import (
	"fmt"
	"io"
	"os"

	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/stdprinter"
	"github.com/friedenberg/zit/charlie/line_format"
	"github.com/friedenberg/zit/charlie/open_file_guard"
	"github.com/friedenberg/zit/echo/zettel"
)

func (f Text) WriteTo(c zettel.FormatContextWrite) (n int64, err error) {
	if c.IncludeAkte {
		if c.ExternalAktePath == "" {
			return f.writeToInlineAkte(c)
		} else {
			return f.writeToExternalAkte(c)
		}
	} else {
		return f.writeToOmitAkte(c)
	}
}

func (f Text) writeToOmitAkte(c zettel.FormatContextWrite) (n int64, err error) {
	w := line_format.NewWriter()

	w.WriteLines(
		MetadateiBoundary,
	)

	if c.Zettel.Bezeichnung.String() != "" || !f.DoNotWriteEmptyBezeichnung {
		w.WriteLines(
			fmt.Sprintf("# %s", c.Zettel.Bezeichnung),
		)
	}

	for _, e := range c.Zettel.Etiketten.Sorted() {
		w.WriteFormat("- %s", e)
	}

	if c.Zettel.Akte.IsNull() && c.Zettel.AkteExt.String() == "" {
		//no-op
	} else if c.Zettel.Akte.IsNull() {
		w.WriteLines(
			fmt.Sprintf("! %s", c.Zettel.AkteExt),
		)
	} else if c.Zettel.AkteExt.String() == "" {
		w.WriteLines(
			fmt.Sprintf("! %s", c.Zettel.Akte),
		)
	} else {
		w.WriteLines(
			fmt.Sprintf("! %s.%s", c.Zettel.Akte, c.Zettel.AkteExt),
		)
	}

	w.WriteLines(
		MetadateiBoundary,
	)

	n, err = w.WriteTo(c.Out)

	return
}

func (f Text) writeToInlineAkte(c zettel.FormatContextWrite) (n int64, err error) {
	if c.Out == nil {
		err = errors.Errorf("context.Out is empty")
		return
	}

	w := line_format.NewWriter()

	w.WriteLines(
		MetadateiBoundary,
		fmt.Sprintf("# %s", c.Zettel.Bezeichnung),
	)

	for _, e := range c.Zettel.Etiketten.Sorted() {
		w.WriteFormat("- %s", e)
	}

	w.WriteLines(
		fmt.Sprintf("! %s", c.Zettel.AkteExt),
	)

	w.WriteLines(
		MetadateiBoundary,
	)

	w.WriteEmpty()

	n, err = w.WriteTo(c.Out)

	if err != nil {
		err = errors.Error(err)
		return
	}

	var ar io.ReadCloser

	if c.AkteReaderFactory == nil {
		err = errors.Errorf("akte reader factory is nil")
		return
	}

	ar, err = c.AkteReader(c.Zettel.Akte)

	if err != nil {
		err = errors.Error(err)
		return
	}

	if ar == nil {
		err = errors.Errorf("akte reader is nil")
		return
	}

	defer stdprinter.PanicIfError(ar.Close)

	var n1 int64

	n1, err = io.Copy(c.Out, ar)
	n += n1

	if err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (f Text) writeToExternalAkte(c zettel.FormatContextWrite) (n int64, err error) {
	logz.Print()
	w := line_format.NewWriter()

	w.WriteLines(
		MetadateiBoundary,
		fmt.Sprintf("# %s", c.Zettel.Bezeichnung),
	)

	for _, e := range c.Zettel.Etiketten.Sorted() {
		w.WriteFormat("- %s", e)
	}

	w.WriteLines(
		fmt.Sprintf("! %s", c.ExternalAktePath),
	)

	w.WriteLines(
		MetadateiBoundary,
	)

	n, err = w.WriteTo(c.Out)

	if err != nil {
		err = errors.Error(err)
		return
	}

	var ar io.ReadCloser

	if c.AkteReaderFactory == nil {
		err = errors.Errorf("akte reader factory is nil")
		return
	}

	if ar, err = c.AkteReader(c.Zettel.Akte); err != nil {
		err = errors.Error(err)
		return
	}

	if ar == nil {
		err = errors.Errorf("akte reader is nil")
		return
	}

	defer stdprinter.PanicIfError(ar.Close)

	var file *os.File

	if file, err = open_file_guard.Create(c.ExternalAktePath); err != nil {
		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(file)

	var n1 int64

	n1, err = io.Copy(file, ar)
	n += n1

	if err != nil {
		err = errors.Error(err)
		return
	}

	return
}
