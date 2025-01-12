package repo_layout

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
)

type ObjectStore struct {
	Config
	basePath string
	age      *age.Age
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
		Age:             s.age,
		Path:            id.Path(sh.GetShaLike(), p),
		CompressionType: s.compressionType,
	}

	if rc, err = dir_layout.NewFileReader(o); err != nil {
		err = errors.Wrapf(err, "Genre: %s", g.GetGenre())
		err = errors.Wrapf(err, "Sha: %s", sh.GetShaLike())
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
		Age:                      s.age,
		FinalPath:                p,
		GenerateFinalPathFromSha: true,
		LockFile:                 s.lockInternalFiles,
		CompressionType:          s.compressionType,
		TemporaryFS:              s.TemporaryFS,
	}

	if wc, err = dir_layout.NewMover(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
