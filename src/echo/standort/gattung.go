package standort

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

func (s Standort) DirObjektenGattung(
	g gattung.GattungLike,
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

	default:
		err = errors.Errorf("unsupported gattung: %s", g)
		return
	}

	return
}

func (s Standort) DirObjektenKennungen() string {
	return s.DirObjekten("Kennungen")
}

func (s Standort) DirObjektenZettelen() string {
	return s.DirObjekten("Zettelen")
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
