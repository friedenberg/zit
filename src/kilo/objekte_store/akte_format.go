package objekte_store

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

type akteFormat[
	O schnittstellen.Akte[O],
	OPtr schnittstellen.AktePtr[O],
] struct {
	objekte.AkteParser[O, OPtr]
	objekte.AkteParseSaver[O, OPtr]
	objekte.SavedAkteFormatter
	objekte.ParsedAkteFormatter[O, OPtr]
}

func MakeAkteFormat[
	O schnittstellen.Akte[O],
	OPtr schnittstellen.AktePtr[O],
](
	akteParser objekte.AkteParser[O, OPtr],
	parsedAkteFormatter objekte.ParsedAkteFormatter[O, OPtr],
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
	e OPtr,
) (n int64, err error) {
	if af.ParsedAkteFormatter == nil {
		err = errors.Errorf("no ParsedAkteFormatter")
	} else {
		n, err = af.ParsedAkteFormatter.FormatParsedAkte(w, e)
	}

	return
}
