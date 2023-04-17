package metadatei

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/format"
)

type FormatterContext interface {
	GetMetadatei() Metadatei
	GetAktePath() string
	GetAkteSha() schnittstellen.Sha
}

type TextFormatter struct {
	DoNotWriteEmptyBezeichnung bool
	IncludeAkteSha             bool
}

func (f *TextFormatter) Format(
	w1 io.Writer,
	c FormatterContext,
) (n int64, err error) {
	w := format.NewLineWriter()
	m := c.GetMetadatei()

	if m.Bezeichnung.String() != "" || !f.DoNotWriteEmptyBezeichnung {
		w.WriteLines(
			fmt.Sprintf("# %s", m.Bezeichnung),
		)
	}

	for _, e := range collections.SortedValues(m.Etiketten) {
		if e.IsEmpty() {
			continue
		}

		w.WriteFormat("- %s", e)
	}

	ap := c.GetAktePath()

	switch {
	case ap != "":
		w.WriteLines(
			fmt.Sprintf("! %s", ap),
		)

	case f.IncludeAkteSha:
		sh := c.GetAkteSha()

		w.WriteLines(
			fmt.Sprintf("! %s.%s", sh, m.Typ),
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
