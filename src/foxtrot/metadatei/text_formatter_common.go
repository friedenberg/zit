package metadatei

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/script_config"
	"github.com/friedenberg/zit/src/delta/ohio"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/echo/format"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type textFormatterCommon struct {
	standort                   standort.Standort
	akteFactory                schnittstellen.AkteReaderFactory
	akteFormatter              script_config.RemoteScript
	doNotWriteEmptyBezeichnung bool
}

func (f textFormatterCommon) writeComments(
	w1 io.Writer,
	c TextFormatterContext,
) (n int64, err error) {
	n1 := 0

	for _, c := range c.GetMetadatei().Comments {
		n1, err = io.WriteString(w1, "% ")
		n += int64(n1)

		if err != nil {
			return
		}

		n1, err = io.WriteString(w1, c)
		n += int64(n1)

		if err != nil {
			return
		}

		n1, err = io.WriteString(w1, "\n")
		n += int64(n1)

		if err != nil {
			return
		}
	}

	return
}

func (f textFormatterCommon) writeBoundary(
	w1 io.Writer,
	_ TextFormatterContext,
) (n int64, err error) {
	return ohio.WriteLine(w1, Boundary)
}

func (f textFormatterCommon) writeNewLine(
	w1 io.Writer,
	_ TextFormatterContext,
) (n int64, err error) {
	return ohio.WriteLine(w1, "")
}

func (f textFormatterCommon) writeCommonMetadateiFormat(
	w1 io.Writer,
	c TextFormatterContext,
) (n int64, err error) {
	w := format.NewLineWriter()
	m := c.GetMetadatei()

	if m.Bezeichnung.String() != "" || !f.doNotWriteEmptyBezeichnung {
		w.WriteLines(
			fmt.Sprintf("# %s", m.Bezeichnung),
		)
	}

	for _, e := range iter.SortedValues[kennung.Etikett](m.GetEtiketten()) {
		if kennung.IsEmpty(e) {
			continue
		}

		w.WriteFormat("- %s", e)
	}

	if n, err = w.WriteTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f textFormatterCommon) writeTyp(
	w1 io.Writer,
	c TextFormatterContext,
) (n int64, err error) {
	m := c.GetMetadatei()
	return ohio.WriteLine(w1, fmt.Sprintf("! %s", m.Typ))
}

func (f textFormatterCommon) writeShaTyp(
	w1 io.Writer,
	c TextFormatterContext,
) (n int64, err error) {
	m := c.GetMetadatei()
	return ohio.WriteLine(w1, fmt.Sprintf("! %s.%s", &m.AkteSha, m.Typ))
}

func (f textFormatterCommon) writePathTyp(
	w1 io.Writer,
	c TextFormatterContext,
) (n int64, err error) {
	var ap string

	if apg, ok := c.(AktePathGetter); ok {
		ap = apg.GetAktePath()
	} else {
		err = errors.Errorf("unable to convert %T int %T", c, apg)
		return
	}

	return ohio.WriteLine(w1, fmt.Sprintf("! %s", ap))
}

func (f textFormatterCommon) writeAkte(
	w1 io.Writer,
	c TextFormatterContext,
) (n int64, err error) {
	var ar io.ReadCloser
	m := c.GetMetadatei()

	if ar, err = f.akteFactory.AkteReader(&m.AkteSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	if ar == nil {
		err = errors.Errorf("akte reader is nil")
		return
	}

	defer errors.DeferredCloser(&err, ar)

	if f.akteFormatter != nil {
		var wt io.WriterTo

		if wt, err = script_config.MakeWriterToWithStdin(
			f.akteFormatter,
			map[string]string{
				"ZIT_BIN": f.standort.Executable(),
			},
			ar,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if n, err = wt.WriteTo(w1); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if n, err = io.Copy(w1, ar); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
