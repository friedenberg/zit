package zettel

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/foxtrot/metadatei_io"
)

const MetadateiBoundary = metadatei_io.Boundary

// TODO switch to three different formats
// metadatei, zettel-akte-external, zettel-akte-inline
func (f Text) WriteTo(c FormatContextWrite) (n int64, err error) {
	switch {
	case c.IncludeAkte && c.ExternalAktePath == "":
		return f.writeToInlineAkte(c)

	case c.IncludeAkte:
		return f.writeToExternalAkte(c)

	default:
		return f.writeToOmitAkte(c)
	}
}

func (f Text) writeToOmitAkte(c FormatContextWrite) (n int64, err error) {
	w := format.NewWriter()

	w.WriteLines(
		MetadateiBoundary,
	)

	if c.Zettel.Bezeichnung.String() != "" || !f.DoNotWriteEmptyBezeichnung {
		w.WriteLines(
			fmt.Sprintf("# %s", c.Zettel.Bezeichnung),
		)
	}

	for _, e := range c.Zettel.Etiketten.Sorted() {
		if e.IsEmpty() {
			continue
		}

		w.WriteFormat("- %s", e)
	}

	switch {
	//TODO log this state
	case c.Zettel.Akte.IsNull() && c.Zettel.Typ.String() == "":
		break

	case c.Zettel.Akte.IsNull():
		w.WriteLines(
			fmt.Sprintf("! %s", c.Zettel.Typ),
		)

	case c.Zettel.Typ.String() == "":
		w.WriteLines(
			fmt.Sprintf("! %s", c.Zettel.Akte),
		)

	default:
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

func (f Text) writeToInlineAkte(c FormatContextWrite) (n int64, err error) {
	if c.Out == nil {
		err = errors.Errorf("context.Out is empty")
		return
	}

	w := format.NewWriter()

	w.WriteLines(
		MetadateiBoundary,
		fmt.Sprintf("# %s", c.Zettel.Bezeichnung),
	)

	for _, e := range c.Zettel.Etiketten.Sorted() {
		w.WriteFormat("- %s", e)
	}

	w.WriteLines(
		fmt.Sprintf("! %s", c.Zettel.Typ),
	)

	w.WriteLines(
		MetadateiBoundary,
	)

	w.WriteEmpty()

	n, err = w.WriteTo(c.Out)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	var ar io.ReadCloser

	if f.AkteFactory == nil {
		err = errors.Errorf("akte reader factory is nil")
		return
	}

	ar, err = f.AkteFactory.AkteReader(c.Zettel.Akte)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if ar == nil {
		err = errors.Errorf("akte reader is nil")
		return
	}

	defer errors.Deferred(&err, ar.Close)

	in := ar

	var cmd *exec.Cmd

	if c.FormatScript != nil {
		if cmd, err = c.FormatScript.Cmd(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if cmd != nil {
		cmd.Stdin = ar
		cmd.Stderr = os.Stderr

		if in, err = cmd.StdoutPipe(); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = cmd.Start(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	var n1 int64

	n1, err = io.Copy(c.Out, in)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if cmd != nil {
		if err = cmd.Wait(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f Text) writeToExternalAkte(c FormatContextWrite) (n int64, err error) {
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

	n, err = w.WriteTo(c.Out)

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
