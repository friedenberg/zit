package env_repo

import (
	"bytes"
	"io"
	"iter"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
)

type blobStore struct {
	env_dir.Config
	basePath string
	tempFS   env_dir.TemporaryFS
}

func MakeBlobStoreFromLayout(
	s Env,
) (bs blobStore, err error) {
	bs = blobStore{
		Config: env_dir.MakeConfigFromImmutableBlobConfig(
			s.GetConfig().ImmutableConfig.GetBlobStoreConfigImmutable(),
		),
		tempFS: s.GetTempLocal(),
	}

	bs.basePath = s.DirBlobs()

	return
}

func MakeBlobStore(
	basePath string,
	config env_dir.Config,
	tempFS env_dir.TemporaryFS,
) blobStore {
	return blobStore{
		Config:   config,
		basePath: basePath,
		tempFS:   tempFS,
	}
}

func (s blobStore) GetBlobStore() interfaces.BlobStore {
	return s
}

func (s blobStore) GetLocalBlobStore() interfaces.LocalBlobStore {
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

func (s blobStore) AllBlobs() iter.Seq2[interfaces.Sha, error] {
	return func(yield func(interfaces.Sha, error) bool) {
		var sh sha.Sha

		for path, err := range files.DirNamesLevel2(s.basePath) {
			if err != nil {
				if !yield(nil, err) {
					return
				}
			}

			if err = sh.SetFromPath(path); err != nil {
				err = errors.Wrap(err)
				if !yield(nil, err) {
					return
				}
			}

			if !yield(&sh, nil) {
				return
			}
		}
	}
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
		if !env_dir.IsErrBlobMissing(err) {
			err = errors.Wrap(err)
		}

		return
	}

	return
}

func (s blobStore) blobWriterTo(p string) (w sha.WriteCloser, err error) {
	mo := env_dir.MoveOptions{
		Config:                   s.Config,
		FinalPath:                p,
		GenerateFinalPathFromSha: true,
		TemporaryFS:              s.tempFS,
	}

	if w, err = env_dir.NewMover(mo); err != nil {
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

	o := env_dir.FileReadOptions{
		Config: s.Config,
		Path:   p,
	}

	if r, err = env_dir.NewFileReader(o); err != nil {
		if errors.IsNotExist(err) {
			shCopy := sha.GetPool().Get()
			shCopy.ResetWithShaLike(sh.GetShaLike())

			err = env_dir.ErrBlobMissing{
				ShaGetter: shCopy,
				Path:      p,
			}
		} else {
			err = errors.Wrapf(
				err,
				"Path: %q, Compression: %q",
				p,
				s.GetBlobCompression(),
			)
		}

		return
	}

	return
}

func MakeCopyingBlobStore(
	env env_local.Env,
	local interfaces.LocalBlobStore,
	remote interfaces.BlobStore,
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
	env_local.Env
	local  interfaces.LocalBlobStore
	remote interfaces.BlobStore
}

func (s CopyingBlobStore) GetBlobStore() interfaces.BlobStore {
	return s
}

func (s CopyingBlobStore) GetLocalBlobStore() interfaces.LocalBlobStore {
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

func (s CopyingBlobStore) AllBlobs() iter.Seq2[interfaces.Sha, error] {
	return s.local.AllBlobs()
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

	if n, err = CopyBlob(s, s.local, s.remote, sh.GetShaLike()); err != nil {
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
	env env_ui.Env,
	dst interfaces.BlobStore,
	src interfaces.BlobStore,
	blobShaGetter interfaces.ShaGetter,
) (n int64, err error) {
	if src == nil {
		return
	}

	blobSha := blobShaGetter.GetShaLike()

	if dst.HasBlob(blobSha) || blobSha.IsNull() {
		err = env_dir.MakeErrAlreadyExists(
			blobSha,
			"",
		)

		return
	}

	return CopyBlob(env, dst, src, blobSha)
}

// TODO make this honor context closure and abort early
func CopyBlob(
	env env_ui.Env,
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
