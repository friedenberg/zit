package dir_layout

import (
	"bytes"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type BlobStore interface {
	HasBlob(sh sha.ShaLike) (ok bool)
	BlobWriter() (w sha.WriteCloser, err error)
	BlobReader(sh sha.ShaLike) (r sha.ReadCloser, err error)
}

type blobStore struct {
	basePath         string
	tempPath         string
	age              *age.Age
	immutable_config immutable_config.Config
	TemporaryFS
}

func MakeBlobStoreFromHome(s DirLayout) (bs blobStore, err error) {
	bs = blobStore{
		age:              s.age,
		immutable_config: s.immutable_config,
		TemporaryFS:      s.TempLocal,
	}

	if bs.basePath, err = s.DirObjectGenre(genres.Blob); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeBlobStore(
	basePath string,
	age *age.Age,
	compressionType immutable_config.CompressionType,
) blobStore {
	return blobStore{
		basePath: basePath,
		age:      age,
		immutable_config: immutable_config.Config{
			CompressionType: compressionType,
		},
	}
}

func (s blobStore) HasBlob(
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

func (s blobStore) BlobWriter() (w sha.WriteCloser, err error) {
	if w, err = s.blobWriterTo(s.basePath); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s blobStore) BlobReader(sh sha.ShaLike) (r sha.ReadCloser, err error) {
	if sh.GetShaLike().IsNull() {
		r = sha.MakeNopReadCloser(io.NopCloser(bytes.NewReader(nil)))
		return
	}

	if r, err = s.blobReaderFrom(sh, s.basePath); err != nil {
		if !IsErrBlobMissing(err) {
			err = errors.Wrap(err)
		}

		return
	}

	return
}

func (s blobStore) blobWriterTo(p string) (w sha.WriteCloser, err error) {
	mo := MoveOptions{
		Age:                      s.age,
		FinalPath:                p,
		GenerateFinalPathFromSha: true,
		LockFile:                 s.immutable_config.LockInternalFiles,
		CompressionType:          s.immutable_config.CompressionType,
		TemporaryFS:              s.TemporaryFS,
	}

	if w, err = NewMover(mo); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s blobStore) blobReaderFrom(
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

func MakeCopyingBlobStore(local, remote BlobStore) CopyingBlobStore {
	if local == nil {
		panic("nil local blob store")
	}

	return CopyingBlobStore{
		local:  local,
		remote: remote,
	}
}

type CopyingBlobStore struct {
	local, remote BlobStore
}

func (s CopyingBlobStore) HasBlob(sh sha.ShaLike) bool {
	if s.local.HasBlob(sh) {
		return true
	}

	if s.remote != nil && s.remote.HasBlob(sh) {
		return true
	}

	return false
}

func (s CopyingBlobStore) BlobWriter() (w sha.WriteCloser, err error) {
	return s.local.BlobWriter()
}

func (s CopyingBlobStore) BlobReader(
	sh sha.ShaLike,
) (r sha.ReadCloser, err error) {
	if s.local.HasBlob(sh) || s.remote == nil {
		return s.local.BlobReader(sh)
	}

	var n int64

	if n, err = CopyBlob(s.local, s.remote, sh.GetShaLike()); err != nil {
		err = errors.Wrap(err)
		return
	}

	ui.Err().Printf("copied Blob %s (%d bytes)", sh, n)

	return s.local.BlobReader(sh)
}

func CopyBlobIfNecessary(
	dst BlobStore,
	src BlobStore,
	blobShaGetter interfaces.ShaGetter,
) (n int64, err error) {
  if src == nil {
    return
  }

	blobSha := blobShaGetter.GetShaLike()

	if dst.HasBlob(blobSha) || blobSha.IsNull() {
		err = MakeErrAlreadyExists(
			blobSha,
			"",
		)

		return
	}

	return CopyBlob(dst, src, blobSha)
}

func CopyBlob(
	dst BlobStore,
	src BlobStore,
	blobSha interfaces.Sha,
) (n int64, err error) {
  if src == nil {
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
