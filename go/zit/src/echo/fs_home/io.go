package fs_home

import (
	"bytes"
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

func (s Standort) objekteReader(
	g interfaces.GenreGetter,
	sh sha.ShaLike,
) (rc sha.ReadCloser, err error) {
	var p string

	if p, err = s.DirObjektenGattung(
		s.immutable_config.GetStoreVersion(),
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

func (s Standort) objekteWriter(
	g interfaces.GenreGetter,
) (wc sha.WriteCloser, err error) {
	var p string

	if p, err = s.DirObjektenGattung(
		s.immutable_config.GetStoreVersion(),
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

func (s Standort) ReadCloserObjekten(p string) (sha.ReadCloser, error) {
	o := FileReadOptions{
		Age:             s.age,
		Path:            p,
		CompressionType: s.immutable_config.CompressionType,
	}

	return NewFileReader(o)
}

func (s Standort) ReadCloserCache(p string) (sha.ReadCloser, error) {
	o := FileReadOptions{
		Age:             s.age,
		Path:            p,
		CompressionType: s.immutable_config.CompressionType,
	}

	return NewFileReader(o)
}

func (s Standort) WriteCloserObjekten(p string) (w sha.WriteCloser, err error) {
	return s.NewMover(
		MoveOptions{
			Age:             s.age,
			FinalPath:       p,
			LockFile:        s.immutable_config.LockInternalFiles,
			CompressionType: s.immutable_config.CompressionType,
		},
	)
}

func (s Standort) WriteCloserCache(
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

func (s Standort) AkteWriterTo(p string) (w sha.WriteCloser, err error) {
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

func (s Standort) BlobWriterToLight(p string) (w sha.WriteCloser, err error) {
	var outer Writer

	mo := MoveOptions{
		Age:                      s.age,
		FinalPath:                p,
		GenerateFinalPathFromSha: true,
		LockFile:                 s.immutable_config.LockInternalFiles,
		CompressionType:          s.immutable_config.CompressionType,
	}

	if outer, err = s.NewMoverLight(mo); err != nil {
		err = errors.Wrap(err)
		return
	}

	w = outer

	return
}

func (s Standort) BlobWriter() (w sha.WriteCloser, err error) {
	var p string

	if p, err = s.DirObjektenGattung(
		s.immutable_config.GetStoreVersion(),
		genres.Blob,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if w, err = s.AkteWriterTo(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Standort) AkteReaderFile(sh sha.ShaLike) (f *os.File, err error) {
	if sh.GetShaLike().IsNull() {
		err = errors.Errorf("sha is null")
		return
	}

	var p string

	if p, err = s.DirObjektenGattung(
		s.immutable_config.GetStoreVersion(),
		genres.Blob,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	p = id.Path(sh.GetShaLike(), p)

	if f, err = files.OpenFile(
		p,
		os.O_RDONLY,
		0o666,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Standort) BlobReader(sh sha.ShaLike) (r sha.ReadCloser, err error) {
	if sh.GetShaLike().IsNull() {
		r = sha.MakeNopReadCloser(io.NopCloser(bytes.NewReader(nil)))
		return
	}

	var p string

	if p, err = s.DirObjektenGattung(
		s.immutable_config.GetStoreVersion(),
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

func (s Standort) BlobReaderFrom(
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
