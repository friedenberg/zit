package typed_blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
)

type RepoStore struct {
	envRepo env_repo.Env
}

func MakeRepoStore(
	dirLayout env_repo.Env,
) RepoStore {
	return RepoStore{
		envRepo: dirLayout,
	}
}

func (a RepoStore) ReadTypedBlob(
	tipe ids.Type,
	blobSha interfaces.Sha,
) (common repo_blobs.Blob, n int64, err error) {
	var reader interfaces.ShaReadCloser

	if reader, err = a.envRepo.BlobReader(blobSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, reader)

	blob := repo_blobs.TypeWithBlob{
		Type: &tipe,
	}

	if n, err = repo_blobs.Coder.DecodeFrom(&blob, reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	common = *blob.Struct

	return
}

func (store RepoStore) WriteTypedBlob(
	tipe ids.Type,
	blob repo_blobs.Blob,
) (sh interfaces.Sha, n int64, err error) {
	var writer interfaces.ShaWriteCloser

	if writer, err = store.envRepo.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, writer)

	if n, err = repo_blobs.Coder.EncodeTo(
		&repo_blobs.TypeWithBlob{
			Type:   &tipe,
			Struct: &blob,
		},
		writer,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = writer.GetShaLike()

	return
}
