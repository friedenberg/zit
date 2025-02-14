package typed_blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
)

type RepoStore struct {
	envRepo             env_repo.Env
	v0                  TypedStore[repo_blobs.V0, *repo_blobs.V0]
	toml_relay_local_v0 TypedStore[repo_blobs.TomlLocalPathV0, *repo_blobs.TomlLocalPathV0]
}

func MakeRepoStore(
	dirLayout env_repo.Env,
) RepoStore {
	return RepoStore{
		envRepo: dirLayout,
		v0: MakeBlobStore(
			dirLayout,
			MakeBlobFormat(
				MakeTomlDecoderIgnoreTomlErrors[repo_blobs.V0](
					dirLayout,
				),
				TomlBlobEncoder[repo_blobs.V0, *repo_blobs.V0]{},
				dirLayout,
			),
			func(a *repo_blobs.V0) {
				a.Reset()
			},
		),
		toml_relay_local_v0: MakeBlobStore(
			dirLayout,
			MakeBlobFormat(
				MakeTomlDecoderIgnoreTomlErrors[repo_blobs.TomlLocalPathV0](
					dirLayout,
				),
				TomlBlobEncoder[repo_blobs.TomlLocalPathV0, *repo_blobs.TomlLocalPathV0]{},
				dirLayout,
			),
			func(a *repo_blobs.TomlLocalPathV0) {
				a.Reset()
			},
		),
	}
}

// func (a RepoStore) GetCommonStore() interfaces.TypedBlobStore[repo_blobs.Blob] {
// 	return a
// }

func (a RepoStore) ReadTypedBlob(
	tipe interfaces.ObjectId,
	blobSha interfaces.Sha,
) (common repo_blobs.Blob, n int64, err error) {
	var reader interfaces.ShaReadCloser

	if reader, err = a.envRepo.BlobReader(blobSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, reader)

	var blob repo_blobs.TypeWithBlob

	if n, err = repo_blobs.Coder.DecodeFrom(&blob, reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	common = *blob.Object

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
			Object: &blob,
		},
		writer,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = writer.GetShaLike()

	return
}
