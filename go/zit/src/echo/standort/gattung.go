package standort

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

func (s Standort) DirObjektenGattung(
	sv interfaces.StoreVersion,
	g interfaces.GattungGetter,
) (p string, err error) {
	switch sv.GetInt() {
	case 0, 1:
		return s.dirObjektenGattung(g)

	default:
		return s.dirObjektenGattung2(g)
	}
}

func (s Standort) dirObjektenGattung2(
	g1 interfaces.GattungGetter,
) (p string, err error) {
	g := g1.GetGattung()

	if g == gattung.Unknown {
		err = gattung.MakeErrUnsupportedGattung(g)
		return
	}

	p = s.DirObjekten2(g.GetGattungStringPlural())

	return
}

func (s Standort) dirObjektenGattung(
	g interfaces.GattungGetter,
) (p string, err error) {
	switch g.GetGattung() {
	case gattung.Akte:
		p = s.DirObjekten("Akten")

	case gattung.Konfig:
		p = s.DirObjekten("Konfig")

	case gattung.Etikett:
		p = s.DirObjekten("Etiketten")

	case gattung.Typ:
		p = s.DirObjekten("Typen")

	case gattung.Zettel:
		p = s.DirObjekten("Zettelen")

	case gattung.Bestandsaufnahme:
		p = s.DirObjekten("Bestandsaufnahme")

	case gattung.Kasten:
		p = s.DirObjekten("Kasten")

	default:
		err = gattung.MakeErrUnsupportedGattung(g)
		return
	}

	return
}

func (s Standort) HasObjekte(
	sv interfaces.StoreVersion,
	g interfaces.GattungGetter,
	sh sha.ShaLike,
) (ok bool) {
	var d string
	var err error

	if d, err = s.DirObjektenGattung(sv, g); err != nil {
		return
	}

	p := id.Path(sh.GetShaLike(), d)
	ok = files.Exists(p)

	return
}

func (s Standort) HasAkte(
	sv interfaces.StoreVersion,
	sh sha.ShaLike,
) (ok bool) {
	if sh.GetShaLike().IsNull() {
		ok = true
		return
	}

	var d string
	var err error

	if d, err = s.DirObjektenGattung(sv, gattung.Akte); err != nil {
		return
	}

	p := id.Path(sh.GetShaLike(), d)
	ok = files.Exists(p)

	return
}

func (s Standort) DirObjektenTransaktion() string {
	return s.DirObjekten("Transaktion")
}

func (s Standort) ReadAllLevel2Files(
	p string,
	w interfaces.FuncIter[string],
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
	w interfaces.FuncIter[*sha.Sha],
) (err error) {
	wf := func(p string) (err error) {
		var sh *sha.Sha

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
	sv interfaces.StoreVersion,
	g interfaces.GattungGetter,
	w interfaces.FuncIter[*sha.Sha],
) (err error) {
	var p string

	if p, err = s.DirObjektenGattung(sv, g); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.ReadAllShas(p, w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
