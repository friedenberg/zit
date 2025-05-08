package blob_store

import (
	"io"
	"iter"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
)

func MakeCopyingBlobStore(
	env env_ui.Env,
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
	env_ui.Env
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

	if n, err = CopyBlob(s, s.local, s.remote, sh.GetShaLike(), nil); err != nil {
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
	extraWriter io.Writer,
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

	return CopyBlob(env, dst, src, blobSha, extraWriter)
}

// TODO make this honor context closure and abort early
func CopyBlob(
	env env_ui.Env,
	dst interfaces.BlobStore,
	src interfaces.BlobStore,
	blobSha interfaces.Sha,
	extraWriter io.Writer,
) (n int64, err error) {
	if src == nil {
		return
	}

	var rc sha.ReadCloser

	if rc, err = src.BlobReader(blobSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer env.MustClose(rc)

	var wc sha.WriteCloser

	if wc, err = dst.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO should this be closed with an error when the shas don't match to
	// prevent a garbage object in the store?
	defer env.MustClose(wc)

	outputWriter := io.Writer(wc)

	if extraWriter != nil {
		outputWriter = io.MultiWriter(outputWriter, extraWriter)
	}

	if n, err = io.Copy(outputWriter, rc); err != nil {
		err = errors.Wrap(err)
		return
	}

	shaRc := rc.GetShaLike()
	shaWc := wc.GetShaLike()

	if !shaRc.EqualsSha(blobSha) || !shaWc.EqualsSha(blobSha) {
		err = errors.ErrorWithStackf(
			"lookup sha was %s, read sha was %s, but written sha was %s",
			blobSha,
			shaRc,
			shaWc,
		)
	}

	return
}
