package blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/tag_blobs"
	"code.linenisgreat.com/zit/go/zit/src/delta/type_blobs"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_blobs"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/mutable_config"
)

type Store[
	A interfaces.Blob[A],
	APtr interfaces.BlobPtr[A],
] interface {
	SaveBlobText(APtr) (interfaces.Sha, int64, error)
	Format[A, APtr]
	interfaces.BlobGetterPutter[APtr]
}

type VersionedStores struct {
	tag_v0    Store[tag_blobs.V0, *tag_blobs.V0]
	tag_v1    Store[tag_blobs.V1, *tag_blobs.V1]
	repo_v0   Store[repo_blobs.V0, *repo_blobs.V0]
	config_v0 Store[mutable_config.Blob, *mutable_config.Blob]
	type_v0   Store[type_blobs.V0, *type_blobs.V0]
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
				MakeTextParserIgnoreTomlErrors[mutable_config.Blob](
					st,
				),
				ParsedBlobTomlFormatter[mutable_config.Blob, *mutable_config.Blob]{},
				st,
			),
			func(a *mutable_config.Blob) {
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

func (a *VersionedStores) GetConfigV0() Store[mutable_config.Blob, *mutable_config.Blob] {
	return a.config_v0
}

func (a *VersionedStores) GetTypeV0() Store[type_blobs.V0, *type_blobs.V0] {
	return a.type_v0
}
