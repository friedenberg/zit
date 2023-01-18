package standort

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/schnittstellen"
)

// TODO-P4 switch to gattung-matched directories
func (s Standort) DirObjektenGattung(
	g schnittstellen.GattungGetter,
) (p string, err error) {
	switch g.GetGattung() {
	case gattung.Konfig:
		p = s.DirObjektenKonfig()

	case gattung.Etikett:
		p = s.DirObjektenEtiketten()

	case gattung.Typ:
		p = s.DirObjektenTypen()

	case gattung.Zettel:
		p = s.DirObjektenZettelen()

	case gattung.Bestandsaufnahme:
		p = s.DirObjektenBestandsaufnahme()

	default:
		err = gattung.ErrUnsupportedGattung
		err = errors.Wrapf(err, "Gattung: %s", g)
		return
	}

	return
}

func (s Standort) HasObjekte(g schnittstellen.GattungGetter, sh sha.ShaLike) (ok bool) {
	var d string
	var err error

	if d, err = s.DirObjektenGattung(g); err != nil {
		return
	}

	p := id.Path(sh.GetSha(), d)
	ok = files.Exists(p)

	return
}

func (s Standort) HasAkte(sh sha.ShaLike) (ok bool) {
	var d string
	var err error

	if d, err = s.DirObjektenGattung(gattung.Akte); err != nil {
		return
	}

	p := id.Path(sh.GetSha(), d)
	ok = files.Exists(p)

	return
}

func (s Standort) DirObjektenKennungen() string {
	return s.DirObjekten("Kennungen")
}

func (s Standort) DirObjektenZettelen() string {
	return s.DirObjekten("Zettelen")
}

func (s Standort) DirObjektenBestandsaufnahme() string {
	return s.DirObjekten("Bestandsaufnahme")
}

func (s Standort) DirObjektenKonfig() string {
	return s.DirObjekten("Konfig")
}

func (s Standort) DirObjektenTypen() string {
	return s.DirObjekten("Typen")
}

func (s Standort) DirObjektenEtiketten() string {
	return s.DirObjekten("Etiketten")
}

func (s Standort) DirObjektenTransaktion() string {
	return s.DirObjekten("Transaktion")
}

func (s Standort) DirObjektenAkten() string {
	return s.DirObjekten("Akten")
}
