package blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/tag_blob"
	"code.linenisgreat.com/zit/go/zit/src/delta/type_blob"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_blob"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/mutable_config"
)

type Store[
	A interfaces.Blob[A],
	APtr interfaces.BlobPtr[A],
] interface {
	SaveBlobText(APtr) (interfaces.ShaLike, int64, error)
	Format[A, APtr]
	interfaces.BlobGetterPutter[APtr]
}

type VersionedStores struct {
	tag_v0    Store[tag_blob.V0, *tag_blob.V0]
	tag_v1    Store[tag_blob.V1, *tag_blob.V1]
	repo_v0   Store[repo_blob.V0, *repo_blob.V0]
	config_v0 Store[mutable_config.Blob, *mutable_config.Blob]
	type_v0   Store[type_blob.V0, *type_blob.V0]
}

func Make(
	st standort.Standort,
) *VersionedStores {
	return &VersionedStores{
		tag_v0: MakeBlobStore(
			st,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[tag_blob.V0](
					st,
				),
				ParsedBlobTomlFormatter[tag_blob.V0, *tag_blob.V0]{},
				st,
			),
			func(a *tag_blob.V0) {
				a.Reset()
			},
		),
		tag_v1: MakeBlobStore(
			st,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[tag_blob.V1](
					st,
				),
				ParsedBlobTomlFormatter[tag_blob.V1, *tag_blob.V1]{},
				st,
			),
			func(a *tag_blob.V1) {
				a.Reset()
			},
		),
		repo_v0: MakeBlobStore(
			st,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[repo_blob.V0](
					st,
				),
				ParsedBlobTomlFormatter[repo_blob.V0, *repo_blob.V0]{},
				st,
			),
			func(a *repo_blob.V0) {
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
				MakeTextParserIgnoreTomlErrors[type_blob.V0](
					st,
				),
				ParsedBlobTomlFormatter[type_blob.V0, *type_blob.V0]{},
				st,
			),
			func(a *type_blob.V0) {
				a.Reset()
			},
		),
	}
}

func (a *VersionedStores) GetTagV0() Store[tag_blob.V0, *tag_blob.V0] {
	return a.tag_v0
}

func (a *VersionedStores) GetTagV1() Store[tag_blob.V1, *tag_blob.V1] {
	return a.tag_v1
}

func (a *VersionedStores) GetRepoV0() Store[repo_blob.V0, *repo_blob.V0] {
	return a.repo_v0
}

func (a *VersionedStores) GetConfigV0() Store[mutable_config.Blob, *mutable_config.Blob] {
	return a.config_v0
}

func (a *VersionedStores) GetTypeV0() Store[type_blob.V0, *type_blob.V0] {
	return a.type_v0
}
