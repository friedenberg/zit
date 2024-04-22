package akten

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
)

type akteFormat[
	O schnittstellen.Akte[O],
	OPtr schnittstellen.AktePtr[O],
] struct {
	AkteParser[O, OPtr]
	AkteParseSaver[O, OPtr]
	SavedAkteFormatter
	ParsedAkteFormatter[O, OPtr]
}

func MakeAkteFormat[
	O schnittstellen.Akte[O],
	OPtr schnittstellen.AktePtr[O],
](
	akteParser AkteParser[O, OPtr],
	parsedAkteFormatter ParsedAkteFormatter[O, OPtr],
	arf schnittstellen.AkteReaderFactory,
) AkteFormat[O, OPtr] {
	return akteFormat[O, OPtr]{
		AkteParser:          akteParser,
		ParsedAkteFormatter: parsedAkteFormatter,
		SavedAkteFormatter:  MakeSavedAkteFormatter(arf),
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
