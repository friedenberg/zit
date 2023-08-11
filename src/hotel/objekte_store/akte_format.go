package objekte_store

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type akteFormat[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
] struct {
	objekte.AkteParser[OPtr]
	objekte.AkteParseSaver[OPtr]
	objekte.SavedAkteFormatter
	objekte.ParsedAkteFormatter[O]
}

func MakeAkteFormat[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
](
	akteParser objekte.AkteParser[OPtr],
	parsedAkteFormatter objekte.ParsedAkteFormatter[O],
	arf schnittstellen.AkteReaderFactory,
) objekte.AkteFormat[O, OPtr] {
	return akteFormat[O, OPtr]{
		AkteParser:          akteParser,
		ParsedAkteFormatter: parsedAkteFormatter,
		SavedAkteFormatter:  objekte.MakeSavedAkteFormatter(arf),
	}
}

func (af akteFormat[O, OPtr]) FormatParsedAkte(
	w io.Writer,
	e O,
) (n int64, err error) {
	if af.ParsedAkteFormatter == nil {
		err = errors.Errorf("no ParsedAkteFormatter")
	} else {
		n, err = af.ParsedAkteFormatter.FormatParsedAkte(w, e)
	}

	return
}
