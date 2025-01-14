package repo_layout

import (
	"bytes"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
)

type blobStore struct {
	dir_layout.Config
	basePath string
	tempFS   dir_layout.TemporaryFS
}

func MakeBlobStoreFromHome(
	s Layout,
) (bs blobStore, err error) {
	bs = blobStore{
		Config: s.config.BlobStoreImmutableConfig,
		tempFS: s.TempLocal,
	}

	if bs.basePath, err = s.DirObjectGenre(genres.Blob); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeBlobStore(
	basePath string,
	config dir_layout.Config,
	tempFS dir_layout.TemporaryFS,
) blobStore {
	ui.DebugBatsTestBody().Printf("%#v", config)
	return blobStore{
		Config:   config,
		basePath: basePath,
		tempFS:   tempFS,
	}
}

func (s blobStore) GetBlobStore() interfaces.BlobStore {
	return s
}

func (s blobStore) HasBlob(
	sh interfaces.Sha,
) (ok bool) {
	if sh.GetShaLike().IsNull() {
		ok = true
		return
	}

	p := id.Path(sh.GetShaLike(), s.basePath)
	ok = files.Exists(p)

	return
}

func (s blobStore) BlobWriter() (w interfaces.ShaWriteCloser, err error) {
	if w, err = s.blobWriterTo(s.basePath); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s blobStore) BlobReader(
	sh interfaces.Sha,
) (r interfaces.ShaReadCloser, err error) {
	if sh.GetShaLike().IsNull() {
		r = sha.MakeNopReadCloser(io.NopCloser(bytes.NewReader(nil)))
		return
	}

	if r, err = s.blobReaderFrom(sh, s.basePath); err != nil {
		if !dir_layout.IsErrBlobMissing(err) {
			err = errors.Wrap(err)
		}

		return
	}

	return
}

func (s blobStore) blobWriterTo(p string) (w sha.WriteCloser, err error) {
	mo := dir_layout.MoveOptions{
		Config:                   s.Config,
		FinalPath:                p,
		GenerateFinalPathFromSha: true,
		TemporaryFS:              s.tempFS,
	}

	if w, err = dir_layout.NewMover(mo); err != nil {
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

	o := dir_layout.FileReadOptions{
		Config: s.Config,
		Path:   p,
	}

	if r, err = dir_layout.NewFileReader(o); err != nil {
		if errors.IsNotExist(err) {
			shCopy := sha.GetPool().Get()
			shCopy.ResetWithShaLike(sh.GetShaLike())

			err = dir_layout.ErrBlobMissing{
				ShaGetter: shCopy,
				Path:      p,
			}
		} else {
			err = errors.Wrapf(err, "Path: %q, Compression: %q", p, s.GetCompressionType())
		}

		return
	}

	return
}

func MakeCopyingBlobStore(
	env *env.Env,
	local, remote interfaces.BlobStore,
) CopyingBlobStore {
	if local == nil {
		panic("nil local blob store")
	}

	return CopyingBlobStore{
		Env:    env,
		local:  local,
		remote: remote,
	}
}

type CopyingBlobStore struct {
	*env.Env
	local, remote interfaces.BlobStore
}

func (s CopyingBlobStore) GetBlobStore() interfaces.BlobStore {
	return s
}

func (s CopyingBlobStore) HasBlob(sh interfaces.Sha) bool {
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
	sh interfaces.Sha,
) (r interfaces.ShaReadCloser, err error) {
	if s.local.HasBlob(sh) || s.remote == nil {
		return s.local.BlobReader(sh)
	}

	var n int64

	if n, err = CopyBlob(s.Env, s.local, s.remote, sh.GetShaLike()); err != nil {
		err = errors.Wrap(err)
		return
	}

	ui.Err().Printf("copied Blob %s (%d bytes)", sh, n)

	if r, err = s.local.BlobReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func CopyBlobIfNecessary(
	env *env.Env,
	dst interfaces.BlobStore,
	src interfaces.BlobStore,
	blobShaGetter interfaces.ShaGetter,
) (n int64, err error) {
	if src == nil {
		return
	}

	blobSha := blobShaGetter.GetShaLike()

	if dst.HasBlob(blobSha) || blobSha.IsNull() {
		err = dir_layout.MakeErrAlreadyExists(
			blobSha,
			"",
		)

		return
	}

	return CopyBlob(env, dst, src, blobSha)
}

// TODO make this honor context closure and abort early
func CopyBlob(
	env *env.Env,
	dst interfaces.BlobStore,
	src interfaces.BlobStore,
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

	// TODO should this be closed with an error when the shas don't match to
	// prevent a garbage object in the store?
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
