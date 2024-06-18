package akten

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
)

type akteFormat[
	O schnittstellen.Akte[O],
	OPtr schnittstellen.AktePtr[O],
] struct {
	Parser[O, OPtr]
	ParseSaver[O, OPtr]
	SavedAkteFormatter
	ParsedAkteFormatter[O, OPtr]
}

func MakeAkteFormat[
	O schnittstellen.Akte[O],
	OPtr schnittstellen.AktePtr[O],
](
	akteParser Parser[O, OPtr],
	parsedAkteFormatter ParsedAkteFormatter[O, OPtr],
	arf schnittstellen.AkteReaderFactory,
) Format[O, OPtr] {
	return akteFormat[O, OPtr]{
		Parser:              akteParser,
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
