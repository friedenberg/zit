package blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/tag_blobs"
	"code.linenisgreat.com/zit/go/zit/src/delta/type_blobs"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_blobs"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/mutable_config_blobs"
)

type Store[
	A interfaces.Blob[A],
	APtr interfaces.BlobPtr[A],
] interface {
	SaveBlobText(APtr) (interfaces.Sha, int64, error)
	Format[A, APtr]
	interfaces.BlobGetterPutter[APtr]
}

// TODO switch to interfaces instead of structs
type VersionedStores struct {
	tag_v0       Store[tag_blobs.V0, *tag_blobs.V0]
	tag_v1       Store[tag_blobs.V1, *tag_blobs.V1]
	repo_v0      Store[repo_blobs.V0, *repo_blobs.V0]
	config_v0    Store[mutable_config_blobs.V0, *mutable_config_blobs.V0]
	type_v0      Store[type_blobs.V0, *type_blobs.V0]
	type_toml_v1 Store[type_blobs.TomlV1, *type_blobs.TomlV1]
}

func Make(
	st fs_home.Home,
) *VersionedStores {
	return &VersionedStores{
		tag_v0: MakeBlobStore(
			st,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[tag_blobs.V0](
					st,
				),
				ParsedBlobTomlFormatter[tag_blobs.V0, *tag_blobs.V0]{},
				st,
			),
			func(a *tag_blobs.V0) {
				a.Reset()
			},
		),
		tag_v1: MakeBlobStore(
			st,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[tag_blobs.V1](
					st,
				),
				ParsedBlobTomlFormatter[tag_blobs.V1, *tag_blobs.V1]{},
				st,
			),
			func(a *tag_blobs.V1) {
				a.Reset()
			},
		),
		repo_v0: MakeBlobStore(
			st,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[repo_blobs.V0](
					st,
				),
				ParsedBlobTomlFormatter[repo_blobs.V0, *repo_blobs.V0]{},
				st,
			),
			func(a *repo_blobs.V0) {
				a.Reset()
			},
		),
		config_v0: MakeBlobStore(
			st,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[mutable_config_blobs.V0](
					st,
				),
				ParsedBlobTomlFormatter[mutable_config_blobs.V0, *mutable_config_blobs.V0]{},
				st,
			),
			func(a *mutable_config_blobs.V0) {
				a.Reset()
			},
		),
		type_v0: MakeBlobStore(
			st,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[type_blobs.V0](
					st,
				),
				ParsedBlobTomlFormatter[type_blobs.V0, *type_blobs.V0]{},
				st,
			),
			func(a *type_blobs.V0) {
				a.Reset()
			},
		),
		type_toml_v1: MakeBlobStore(
			st,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[type_blobs.TomlV1](
					st,
				),
				ParsedBlobTomlFormatter[type_blobs.TomlV1, *type_blobs.TomlV1]{},
				st,
			),
			func(a *type_blobs.TomlV1) {
				a.Reset()
			},
		),
	}
}

func (a *VersionedStores) GetTagV0() Store[tag_blobs.V0, *tag_blobs.V0] {
	return a.tag_v0
}

func (a *VersionedStores) GetTagV1() Store[tag_blobs.V1, *tag_blobs.V1] {
	return a.tag_v1
}

func (a *VersionedStores) GetRepoV0() Store[repo_blobs.V0, *repo_blobs.V0] {
	return a.repo_v0
}

func (a *VersionedStores) GetConfigV0() Store[mutable_config_blobs.V0, *mutable_config_blobs.V0] {
	return a.config_v0
}

func (a *VersionedStores) GetTypeV0() Store[type_blobs.V0, *type_blobs.V0] {
	return a.type_v0
}

func (a *VersionedStores) GetTypeV1() Store[type_blobs.TomlV1, *type_blobs.TomlV1] {
	return a.type_toml_v1
}

func (a *VersionedStores) ParseTypeBlob(
	tipe interfaces.ObjectId,
	blobSha interfaces.Sha,
) (common type_blobs.Common, n int64, err error) {
	switch tipe.String() {
	case "", type_blobs.TypeV0:
		store := a.GetTypeV0()
		var blob *type_blobs.V0

		if blob, err = store.GetBlob(blobSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		common = blob

	case type_blobs.TypeV1:
		store := a.GetTypeV1()
		var blob *type_blobs.TomlV1

		if blob, err = store.GetBlob(blobSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		common = blob
	}

	return
}

func (a *VersionedStores) PutTypeBlob(
	tipe interfaces.ObjectId,
	common type_blobs.Common,
) (err error) {
	switch tipe.String() {
	case "", type_blobs.TypeV0:
		if blob, ok := common.(*type_blobs.V0); !ok {
			err = errors.Errorf("expected %T but got %T", blob, common)
			return
		} else {
			a.GetTypeV0().PutBlob(blob)
		}

	case type_blobs.TypeV1:
		if blob, ok := common.(*type_blobs.TomlV1); !ok {
			err = errors.Errorf("expected %T but got %T", blob, common)
			return
		} else {
			a.GetTypeV1().PutBlob(blob)
		}
	}

	return
}
