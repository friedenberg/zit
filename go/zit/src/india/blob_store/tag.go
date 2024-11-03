package blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/tag_blobs"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
)

type TagStore struct {
	tag_toml_v0 Store[tag_blobs.V0, *tag_blobs.V0]
	tag_toml_v1 Store[tag_blobs.TomlV1, *tag_blobs.TomlV1]
}

func MakeTagStore(
	dirLayout dir_layout.DirLayout,
) TagStore {
	return TagStore{
		tag_toml_v0: MakeBlobStore(
			dirLayout,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[tag_blobs.V0](
					dirLayout,
				),
				ParsedBlobTomlFormatter[tag_blobs.V0, *tag_blobs.V0]{},
				dirLayout,
			),
			func(a *tag_blobs.V0) {
				a.Reset()
			},
		),
		tag_toml_v1: MakeBlobStore(
			dirLayout,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[tag_blobs.TomlV1](
					dirLayout,
				),
				ParsedBlobTomlFormatter[tag_blobs.TomlV1, *tag_blobs.TomlV1]{},
				dirLayout,
			),
			func(a *tag_blobs.TomlV1) {
				a.Reset()
			},
		),
	}
}

func (a TagStore) GetCommonStore() CommonStore[tag_blobs.Common] {
	return a
}

func (a TagStore) ParseTypedBlob(
	tipe interfaces.ObjectId,
	blobSha interfaces.Sha,
) (common tag_blobs.Common, n int64, err error) {
	switch tipe.String() {
	case "", builtin_types.TagTypeTomlV0:
		store := a.tag_toml_v0
		var blob *tag_blobs.V0

		if blob, err = store.GetBlob(blobSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		common = blob

	case builtin_types.TagTypeTomlV1:
		store := a.tag_toml_v1
		var blob *tag_blobs.TomlV1

		if blob, err = store.GetBlob(blobSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		common = blob
	}

	return
}

func (a TagStore) PutTypedBlob(
	tipe interfaces.ObjectId,
	common tag_blobs.Common,
) (err error) {
	switch tipe.String() {
	case "", builtin_types.TagTypeTomlV0:
		if blob, ok := common.(*tag_blobs.V0); !ok {
			err = errors.Errorf("expected %T but got %T", blob, common)
			return
		} else {
			a.tag_toml_v0.PutBlob(blob)
		}

	case builtin_types.TagTypeLuaV1:
		if blob, ok := common.(*tag_blobs.TomlV1); !ok {
			err = errors.Errorf("expected %T but got %T", blob, common)
			return
		} else {
			a.tag_toml_v1.PutBlob(blob)
		}
	}

	return
}
