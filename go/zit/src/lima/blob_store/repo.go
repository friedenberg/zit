package blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_blobs"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
)

type RepoStore struct {
	dirLayout           env_repo.Env
	v0                  Store[repo_blobs.V0, *repo_blobs.V0]
	toml_relay_local_v0 Store[repo_blobs.TomlRelayLocalV0, *repo_blobs.TomlRelayLocalV0]
}

func MakeRepoStore(
	dirLayout env_repo.Env,
) RepoStore {
	return RepoStore{
		dirLayout: dirLayout,
		v0: MakeBlobStore(
			dirLayout,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[repo_blobs.V0](
					dirLayout,
				),
				ParsedBlobTomlFormatter[repo_blobs.V0, *repo_blobs.V0]{},
				dirLayout,
			),
			func(a *repo_blobs.V0) {
				a.Reset()
			},
		),
		toml_relay_local_v0: MakeBlobStore(
			dirLayout,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[repo_blobs.TomlRelayLocalV0](
					dirLayout,
				),
				ParsedBlobTomlFormatter[repo_blobs.TomlRelayLocalV0, *repo_blobs.TomlRelayLocalV0]{},
				dirLayout,
			),
			func(a *repo_blobs.TomlRelayLocalV0) {
				a.Reset()
			},
		),
	}
}

func (a RepoStore) GetCommonStore() interfaces.TypedBlobStore[repo_blobs.Blob] {
	return a
}

func (a RepoStore) ParseTypedBlob(
	tipe interfaces.ObjectId,
	blobSha interfaces.Sha,
) (common repo_blobs.Blob, n int64, err error) {
	switch tipe.String() {
	case "":
		store := a.v0
		var blob *repo_blobs.V0

		if blob, err = store.GetBlob(blobSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		common = blob
	case builtin_types.RepoTypeLocalRelay:
		store := a.toml_relay_local_v0
		var blob *repo_blobs.TomlRelayLocalV0

		if blob, err = store.GetBlob(blobSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		common = blob
	}

	return
}

func (a RepoStore) PutTypedBlob(
	tipe interfaces.ObjectId,
	common repo_blobs.Blob,
) (err error) {
	switch tipe.String() {
	case "":
		if blob, ok := common.(*repo_blobs.V0); !ok {
			err = errors.Errorf("expected %T but got %T", blob, common)
			return
		} else {
			a.v0.PutBlob(blob)
		}

	case builtin_types.RepoTypeLocalRelay:
		if blob, ok := common.(*repo_blobs.TomlRelayLocalV0); !ok {
			err = errors.Errorf("expected %T but got %T", blob, common)
			return
		} else {
			a.toml_relay_local_v0.PutBlob(blob)
		}
	}

	return
}
