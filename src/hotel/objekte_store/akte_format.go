package objekte_store

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type AkteFormat[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
] interface {
	objekte.AkteParser[OPtr]
	objekte.AkteParseSaver[OPtr]
	objekte.SavedAkteFormatter
	objekte.ParsedAkteFormatter[O]
}

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
) AkteFormat[O, OPtr] {
	return akteFormat[O, OPtr]{
		AkteParseSaver:      akteParseSaver,
		ParsedAkteFormatter: parsedAkteFormatter,
		SavedAkteFormatter:  objekte.MakeSavedAkteFormatter(arf),
	}
}

// func (af akteFormat[O, OPtr]) Format(w io.Writer, e OPtr) (int64, error) {
// 	return af.FormatSavedAkte(w, e.GetAkteSha())
// }

func (af akteFormat[O, OPtr]) ParseAkte(
	r io.Reader,
	e OPtr,
) (n int64, err error) {
	var sh schnittstellen.Sha

	sh, n, err = af.ParseSaveAkte(r, e)

	e.SetAkteSha(sha.Make(sh))

	return
}

func (af akteFormat[O, OPtr]) FormatParsedAkte(
	w io.Writer,
	e O,
) (n int64, err error) {
	if af.ParsedAkteFormatter == nil {
		sh := e.GetAkteSha()
		n, err = af.SavedAkteFormatter.FormatSavedAkte(w, sh)
	} else {
		n, err = af.ParsedAkteFormatter.FormatParsedAkte(w, e)
	}

	return
}
