package standort

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
)

// TODO-P4 switch to gattung-matched directories
func (s Standort) DirObjektenGattung(
	g schnittstellen.GattungGetter,
) (p string, err error) {
	switch g.GetGattung() {
	case gattung.Akte:
		p = s.DirObjektenAkten()

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

	case gattung.Kasten:
		p = s.DirObjektenKasten()

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

func (s Standort) DirObjektenKasten() string {
	return s.DirObjekten("Kasten")
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

func (s Standort) ReadAllLevel2Files(
	p string,
	w collections.WriterFunc[string],
) (err error) {
	if err = files.ReadDirNamesLevel2(
		files.MakeDirNameWriterIgnoringHidden(w),
		p,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Standort) ReadAllShas(
	p string,
	w collections.WriterFunc[sha.Sha],
) (err error) {
	wf := func(p string) (err error) {
		var sh sha.Sha

		if sh, err = sha.MakeShaFromPath(p); err != nil {
			err = errors.Wrapf(err, "Path: %s", p)
			return
		}

		if err = w(sh); err != nil {
			err = errors.Wrapf(err, "Sha: %s", sh)
			return
		}

		return
	}

	if err = s.ReadAllLevel2Files(p, wf); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Standort) ReadAllShasForGattung(
	g schnittstellen.GattungGetter,
	w collections.WriterFunc[sha.Sha],
) (err error) {
	var p string

	if p, err = s.DirObjektenGattung(g); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.ReadAllShas(p, w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
