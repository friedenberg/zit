package zettel_formats

import (
	"fmt"
	"io"
	"log"
	"os"
)

func (f Text) WriteTo(c _ZettelFormatContextWrite) (n int64, err error) {
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

func (f Text) writeToOmitAkte(c _ZettelFormatContextWrite) (n int64, err error) {
	log.Print()
	w := _LineFormatNewWriter()

	w.WriteLines(
		MetadateiBoundary,
		fmt.Sprintf("# %s", c.Zettel.Bezeichnung),
	)

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

func (f Text) writeToInlineAkte(c _ZettelFormatContextWrite) (n int64, err error) {
	log.Print()
	w := _LineFormatNewWriter()

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
		err = _Error(err)
		return
	}

	var ar io.ReadCloser

	if c.AkteReaderFactory == nil {
		err = _Errorf("akte reader factory is nil")
		return
	}

	ar, err = c.AkteReader(c.Zettel.Akte)

	if err != nil {
		err = _Error(err)
		return
	}

	if ar == nil {
		err = _Errorf("akte reader is nil")
		return
	}

	defer _PanicIfError(ar.Close())

	var n1 int64

	n1, err = io.Copy(c.Out, ar)
	n += n1

	if err != nil {
		err = _Error(err)
		return
	}

	return
}

func (f Text) writeToExternalAkte(c _ZettelFormatContextWrite) (n int64, err error) {
	log.Print()
	w := _LineFormatNewWriter()

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
		err = _Error(err)
		return
	}

	var ar io.ReadCloser

	if c.AkteReaderFactory == nil {
		err = _Errorf("akte reader factory is nil")
		return
	}

	ar, err = c.AkteReader(c.Zettel.Akte)

	if err != nil {
		err = _Error(err)
		return
	}

	if ar == nil {
		err = _Errorf("akte reader is nil")
		return
	}

	defer _PanicIfError(ar.Close())

	var file *os.File

	if file, err = _Create(c.ExternalAktePath); err != nil {
		err = _Error(err)
		return
	}

	defer _Close(file)

	var n1 int64

	n1, err = io.Copy(file, ar)
	n += n1

	if err != nil {
		err = _Error(err)
		return
	}

	return
}
