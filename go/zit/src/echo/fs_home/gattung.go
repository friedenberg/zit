package fs_home

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

func (s Home) DirObjektenGattung(
	sv interfaces.StoreVersion,
	g interfaces.GenreGetter,
) (p string, err error) {
	switch sv.GetInt() {
	case 0, 1:
		return s.dirObjektenGattung(g)

	default:
		return s.dirObjektenGattung2(g)
	}
}

func (s Home) dirObjektenGattung2(
	g1 interfaces.GenreGetter,
) (p string, err error) {
	g := g1.GetGenre()

	if g == genres.Unknown {
		err = genres.MakeErrUnsupportedGenre(g)
		return
	}

	p = s.DirObjekten2(g.GetGenreStringPlural())

	return
}

func (s Home) dirObjektenGattung(
	g interfaces.GenreGetter,
) (p string, err error) {
	switch g.GetGenre() {
	case genres.Blob:
		p = s.DirObjekten("Akten")

	case genres.Config:
		p = s.DirObjekten("Konfig")

	case genres.Tag:
		p = s.DirObjekten("Etiketten")

	case genres.Type:
		p = s.DirObjekten("Typen")

	case genres.Zettel:
		p = s.DirObjekten("Zettelen")

	case genres.InventoryList:
		p = s.DirObjekten("Bestandsaufnahme")

	case genres.Repo:
		p = s.DirObjekten("Kasten")

	default:
		err = genres.MakeErrUnsupportedGenre(g)
		return
	}

	return
}

func (s Home) HasObjekte(
	sv interfaces.StoreVersion,
	g interfaces.GenreGetter,
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

func (s Home) HasBlob(
	sv interfaces.StoreVersion,
	sh sha.ShaLike,
) (ok bool) {
	if sh.GetShaLike().IsNull() {
		ok = true
		return
	}

	var d string
	var err error

	if d, err = s.DirObjektenGattung(sv, genres.Blob); err != nil {
		return
	}

	p := id.Path(sh.GetShaLike(), d)
	ok = files.Exists(p)

	return
}

func (s Home) DirObjektenTransaktion() string {
	return s.DirObjekten("Transaktion")
}

func (s Home) ReadAllLevel2Files(
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

func (s Home) ReadAllShas(
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

func (s Home) ReadAllShasForGenre(
	sv interfaces.StoreVersion,
	g interfaces.GenreGetter,
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
