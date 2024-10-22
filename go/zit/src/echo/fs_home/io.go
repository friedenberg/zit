package fs_home

import (
	"bytes"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

func (s Home) objectReader(
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

func (s Home) objectWriter(
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
	}

	if wc, err = s.NewMover(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Home) ReadCloserObjekten(p string) (sha.ReadCloser, error) {
	o := FileReadOptions{
		Age:             s.age,
		Path:            p,
		CompressionType: s.immutable_config.CompressionType,
	}

	return NewFileReader(o)
}

func (s Home) ReadCloserCache(p string) (sha.ReadCloser, error) {
	o := FileReadOptions{
		Age:             s.age,
		Path:            p,
		CompressionType: s.immutable_config.CompressionType,
	}

	return NewFileReader(o)
}

func (s Home) WriteCloserObjekten(p string) (w sha.WriteCloser, err error) {
	return s.NewMover(
		MoveOptions{
			Age:             s.age,
			FinalPath:       p,
			LockFile:        s.immutable_config.LockInternalFiles,
			CompressionType: s.immutable_config.CompressionType,
		},
	)
}

func (s Home) WriteCloserCache(
	p string,
) (w sha.WriteCloser, err error) {
	return s.NewMover(
		MoveOptions{
			Age:             s.age,
			FinalPath:       p,
			LockFile:        false,
			CompressionType: s.immutable_config.CompressionType,
		},
	)
}

func (s Home) BlobWriterTo(p string) (w sha.WriteCloser, err error) {
	var outer Writer

	mo := MoveOptions{
		Age:                      s.age,
		FinalPath:                p,
		GenerateFinalPathFromSha: true,
		LockFile:                 s.immutable_config.LockInternalFiles,
		CompressionType:          s.immutable_config.CompressionType,
	}

	if outer, err = s.NewMover(mo); err != nil {
		err = errors.Wrap(err)
		return
	}

	w = outer

	return
}

func (s Home) BlobWriter() (w sha.WriteCloser, err error) {
	var p string

	if p, err = s.DirObjectGenre(
		genres.Blob,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if w, err = s.BlobWriterTo(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Home) BlobReader(sh sha.ShaLike) (r sha.ReadCloser, err error) {
	if sh.GetShaLike().IsNull() {
		r = sha.MakeNopReadCloser(io.NopCloser(bytes.NewReader(nil)))
		return
	}

	var p string

	if p, err = s.DirObjectGenre(
		genres.Blob,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if r, err = s.BlobReaderFrom(sh, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Home) BlobReaderFrom(
	sh sha.ShaLike,
	p string,
) (r sha.ReadCloser, err error) {
	if sh.GetShaLike().IsNull() {
		r = sha.MakeNopReadCloser(io.NopCloser(bytes.NewReader(nil)))
		return
	}

	p = id.Path(sh.GetShaLike(), p)

	o := FileReadOptions{
		Age:             s.age,
		Path:            p,
		CompressionType: s.immutable_config.CompressionType,
	}

	if r, err = NewFileReader(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
