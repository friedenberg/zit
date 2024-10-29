package fs_home

import (
	"bytes"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type BlobStore struct {
	basePath         string
	tempPath         string
	age              *age.Age
	immutable_config immutable_config.Config
	TemporaryFS
	MoverFactory
}

func MakeBlobStoreFromHome(s Home) (bs BlobStore, err error) {
	bs = BlobStore{
		age:              s.age,
		immutable_config: s.immutable_config,
		MoverFactory:     s,
		TemporaryFS:      s.TempLocal,
	}

	if bs.basePath, err = s.DirObjectGenre(genres.Blob); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeBlobStore(h Home) BlobStore {
	return BlobStore{
		// basePath:         s.basePath,
		// age:              s.age,
		// immutable_config: s.immutable_config,
		// MoverFactory:     s,
	}
}

func (s BlobStore) HasBlob(
	sh sha.ShaLike,
) (ok bool) {
	if sh.GetShaLike().IsNull() {
		ok = true
		return
	}

	p := id.Path(sh.GetShaLike(), s.basePath)
	ok = files.Exists(p)

	return
}

func (s BlobStore) BlobWriterTo(p string) (w sha.WriteCloser, err error) {
	mo := MoveOptions{
		Age:                      s.age,
		FinalPath:                p,
		GenerateFinalPathFromSha: true,
		LockFile:                 s.immutable_config.LockInternalFiles,
		CompressionType:          s.immutable_config.CompressionType,
	}

	if w, err = s.NewMover(mo); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s BlobStore) BlobWriter() (w sha.WriteCloser, err error) {
	if w, err = s.BlobWriterTo(s.basePath); err != nil {
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

	if r, err = s.BlobReaderFrom(sh, s.basePath); err != nil {
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

func (dst BlobStore) CopyBlobIfNecessary(
	src BlobStore,
	blobShaGetter interfaces.ShaGetter,
) (n int64, err error) {
	blobSha := blobShaGetter.GetShaLike()

	if dst.HasBlob(blobSha) {
		return
	}

	var rc sha.ReadCloser

	if rc, err = src.BlobReader(blobSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, rc)

	var wc sha.WriteCloser

	if wc, err = dst.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, wc)

	if n, err = io.Copy(wc, rc); err != nil {
		err = errors.Wrap(err)
		return
	}

	shaRc := rc.GetShaLike()
	shaWc := wc.GetShaLike()

	if !shaRc.EqualsSha(blobSha) || !shaWc.EqualsSha(blobSha) {
		err = errors.Errorf(
			"lookup sha was %s, read sha was %s, but written sha was %s",
			blobSha,
			shaRc,
			shaWc,
		)
	}

	return
}
