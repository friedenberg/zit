package repo_layout

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
)

type ObjectStore struct {
	dir_layout.Config

	basePath string
	interfaces.DirectoryPaths
	dir_layout.TemporaryFS
}

func (s ObjectStore) objectReader(
	g interfaces.GenreGetter,
	sh sha.ShaLike,
) (rc sha.ReadCloser, err error) {
	var p string

	if p, err = s.DirObjectGenre(
		g,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	o := dir_layout.FileReadOptions{
		Config: s.Config,
		Path:   id.Path(sh.GetShaLike(), p),
	}

	if rc, err = dir_layout.NewFileReader(o); err != nil {
		err = errors.Wrapf(err, "Genre: %s", g.GetGenre())
		err = errors.Wrapf(err, "Sha: %s", sh.GetShaLike())
		err = errors.Wrapf(err, "Path: %s", o.Path)
		err = errors.Wrapf(err, "Age: %s", o.GetAgeEncryption())
		return
	}

	return
}

func (s ObjectStore) objectWriter(
	g interfaces.GenreGetter,
) (wc sha.WriteCloser, err error) {
	var p string

	if p, err = s.DirObjectGenre(
		g,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	o := dir_layout.MoveOptions{
		Config:                   s.Config,
		FinalPath:                p,
		GenerateFinalPathFromSha: true,
		TemporaryFS:              s.TemporaryFS,
	}

	if wc, err = dir_layout.NewMover(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
