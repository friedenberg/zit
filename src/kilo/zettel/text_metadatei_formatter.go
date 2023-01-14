package zettel

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/format"
)

type Metadatei struct {
	Objekte
	AktePath string
}

type TextMetadateiFormatter struct {
	DoNotWriteEmptyBezeichnung bool
	IncludeAkteSha             bool
}

func (f *TextMetadateiFormatter) Format(w1 io.Writer, m *Metadatei) (n int64, err error) {
	errors.Todo(errors.P3, "turn *Objekte into an interface to allow zettel_external to use this")

	w := format.NewWriter()

	if m.Bezeichnung.String() != "" || !f.DoNotWriteEmptyBezeichnung {
		w.WriteLines(
			fmt.Sprintf("# %s", m.Bezeichnung),
		)
	}

	for _, e := range m.Etiketten.Sorted() {
		errors.Todo(errors.P3, "determine how to handle this")

		if e.IsEmpty() {
			continue
		}

		w.WriteFormat("- %s", e)
	}

	switch {
	case m.AktePath != "":
		w.WriteLines(
			fmt.Sprintf("! %s", m.AktePath),
		)

	case f.IncludeAkteSha:
		w.WriteLines(
			fmt.Sprintf("! %s.%s", m.Akte, m.Typ),
		)

	default:
		w.WriteLines(
			fmt.Sprintf("! %s", m.Typ),
		)
	}

	if n, err = w.WriteTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
