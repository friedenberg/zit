package repo_layout

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
)

type ObjectStore struct {
	env_dir.Config

	basePath string
	interfaces.DirectoryPaths
	env_dir.TemporaryFS
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

	o := env_dir.FileReadOptions{
		Config: s.Config,
		Path:   id.Path(sh.GetShaLike(), p),
	}

	if rc, err = env_dir.NewFileReader(o); err != nil {
		err = errors.Wrapf(err, "Genre: %s", g.GetGenre())
		err = errors.Wrapf(err, "Sha: %s", sh.GetShaLike())
		err = errors.Wrapf(err, "Path: %s", o.Path)
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

	o := env_dir.MoveOptions{
		Config:                   s.Config,
		FinalPath:                p,
		GenerateFinalPathFromSha: true,
		TemporaryFS:              s.TemporaryFS,
	}

	if wc, err = env_dir.NewMover(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
