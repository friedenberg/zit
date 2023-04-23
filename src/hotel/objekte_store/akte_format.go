package objekte_store

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type akteFormat[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
] struct {
	objekte.AkteParseSaver[OPtr]
	objekte.SavedAkteFormatter
	objekte.ParsedAkteFormatter[O]
}

func MakeAkteFormat[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
](
	akteParseSaver objekte.AkteParseSaver[OPtr],
	parsedAkteFormatter objekte.ParsedAkteFormatter[O],
	arf schnittstellen.AkteReaderFactory,
) objekte.AkteFormat[O, OPtr] {
	return akteFormat[O, OPtr]{
		AkteParseSaver:      akteParseSaver,
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
