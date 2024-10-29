package fs_home

import (
	"bytes"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type BlobStore struct {
	basePath         string
	age              *age.Age
	immutable_config immutable_config.Config
	interfaces.DirectoryPaths
	MoverFactory
}

func (s BlobStore) BlobWriterTo(p string) (w sha.WriteCloser, err error) {
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

func (s BlobStore) BlobWriter() (w sha.WriteCloser, err error) {
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

func (s BlobStore) BlobReader(sh sha.ShaLike) (r sha.ReadCloser, err error) {
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
		if !IsErrBlobMissing(err) {
			err = errors.Wrap(err)
		}

		return
	}

	return
}

func (s BlobStore) BlobReaderFrom(
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
		if errors.IsNotExist(err) {
			shCopy := sha.GetPool().Get()
			shCopy.ResetWithShaLike(sh.GetShaLike())

			err = ErrBlobMissing{
				ShaGetter: shCopy,
				Path:      p,
			}
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	return
}
