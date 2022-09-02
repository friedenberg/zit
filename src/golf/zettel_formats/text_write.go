package zettel_formats

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/charlie/line_format"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
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

	if c.Zettel.Akte.IsNull() && c.Zettel.Typ.String() == "" {
		//no-op
	} else if c.Zettel.Akte.IsNull() {
		w.WriteLines(
			fmt.Sprintf("! %s", c.Zettel.Typ),
		)
	} else if c.Zettel.Typ.String() == "" {
		w.WriteLines(
			fmt.Sprintf("! %s", c.Zettel.Akte),
		)
	} else {
		w.WriteLines(
			fmt.Sprintf("! %s.%s", c.Zettel.Akte, c.Zettel.Typ),
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
		fmt.Sprintf("! %s", c.Zettel.TypOrDefault()),
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

	in := ar

	var cmd *exec.Cmd

	if c.FormatScript != nil {
		if cmd, err = c.FormatScript.Cmd(); err != nil {
			err = errors.Error(err)
			return
		}

		cmd.Stdin = ar
		cmd.Stderr = os.Stderr

		if in, err = cmd.StdoutPipe(); err != nil {
			err = errors.Error(err)
			return
		}

		if err = cmd.Start(); err != nil {
			err = errors.Error(err)
			return
		}
	}

	var n1 int64

	n1, err = io.Copy(c.Out, in)
	n += n1

	if err != nil {
		err = errors.Error(err)
		return
	}

	if cmd != nil {
		if err = cmd.Wait(); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}

func (f Text) writeToExternalAkte(c zettel.FormatContextWrite) (n int64, err error) {
	errors.Print()
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
