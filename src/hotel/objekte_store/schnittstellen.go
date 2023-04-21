package objekte_store

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type akteFormat[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
] struct {
	objekte.AkteParseSaver[OPtr]
	objekte.AkteFormatter
}

func MakeAkteFormat[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
](
	akteParseSaver objekte.AkteParseSaver[OPtr],
	arf schnittstellen.AkteReaderFactory,
) akteFormat[O, OPtr] {
	return akteFormat[O, OPtr]{
		AkteParseSaver: akteParseSaver,
		AkteFormatter:  objekte.MakeAkteFormatter(arf),
	}
}

func (af akteFormat[O, OPtr]) Format(w io.Writer, e OPtr) (int64, error) {
	return af.FormatAkte(w, e.GetAkteSha())
}

func (af akteFormat[O, OPtr]) Parse(r io.Reader, e OPtr) (n int64, err error) {
	var sh schnittstellen.Sha

	sh, n, err = af.ParseSaveAkte(r, e)

	e.SetAkteSha(sha.Make(sh))

	return
}
