package sku_fmt

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type organize struct {
	maxKopf, maxSchwanz int
	omitNewLine         bool
}

func MakeOrganizeFormat(
	maxKopf, maxSchwanz int,
) *organize {
	return &organize{
		maxKopf:     maxKopf,
		maxSchwanz:  maxSchwanz,
		omitNewLine: true,
	}
}

func (f *organize) WriteStringFormat(
	sw io.StringWriter,
	o *sku.Transacted,
) (n int64, err error) {
	var n1 int

	n1, err = sw.WriteString("- [")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	h := kennung.Aligned(o.Kennung, f.maxKopf, f.maxSchwanz)

	n1, err = sw.WriteString(h)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = sw.WriteString("]")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if !o.Metadatei.Bezeichnung.IsEmpty() {
		n1, err = sw.WriteString(" ")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n1, err = sw.WriteString(o.Metadatei.Description())
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if !f.omitNewLine {
		n1, err = sw.WriteString("\n")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
