package dir_layout

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type ObjectStore struct {
	basePath         string
	age              *age.Age
	immutable_config immutable_config.Config
	interfaces.DirectoryPaths
	TemporaryFS
}

func (s DirLayout) HasObject(
	g interfaces.GenreGetter,
	sh sha.ShaLike,
) (ok bool) {
	var d string
	var err error

	if d, err = s.DirObjectGenre(g); err != nil {
		return
	}

	p := id.Path(sh.GetShaLike(), d)
	ok = files.Exists(p)

	return
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

	o := FileReadOptions{
		Age:             s.age,
		Path:            id.Path(sh.GetShaLike(), p),
		CompressionType: s.immutable_config.CompressionType,
	}

	if rc, err = NewFileReader(o); err != nil {
		err = errors.Wrapf(err, "Gattung: %s", g.GetGenre())
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

	o := MoveOptions{
		Age:                      s.age,
		FinalPath:                p,
		GenerateFinalPathFromSha: true,
		LockFile:                 s.immutable_config.LockInternalFiles,
		CompressionType:          s.immutable_config.CompressionType,
		TemporaryFS:              s.TemporaryFS,
	}

	if wc, err = NewMover(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
