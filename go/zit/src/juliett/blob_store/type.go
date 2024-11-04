package blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/hotel/type_blobs"
)

type TypeStore struct {
	toml_v0 Store[type_blobs.V0, *type_blobs.V0]
	toml_v1 Store[type_blobs.TomlV1, *type_blobs.TomlV1]
}

func MakeTypeStore(
	dirLayout dir_layout.DirLayout,
) TypeStore {
	return TypeStore{
		toml_v0: MakeBlobStore(
			dirLayout,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[type_blobs.V0](
					dirLayout,
				),
				ParsedBlobTomlFormatter[type_blobs.V0, *type_blobs.V0]{},
				dirLayout,
			),
			func(a *type_blobs.V0) {
				a.Reset()
			},
		),
		toml_v1: MakeBlobStore(
			dirLayout,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[type_blobs.TomlV1](
					dirLayout,
				),
				ParsedBlobTomlFormatter[type_blobs.TomlV1, *type_blobs.TomlV1]{},
				dirLayout,
			),
			func(a *type_blobs.TomlV1) {
				a.Reset()
			},
		),
	}
}

func (a TypeStore) GetCommonStore() CommonStore[type_blobs.Blob] {
	return a
}

func (a TypeStore) ParseTypedBlob(
	tipe interfaces.ObjectId,
	blobSha interfaces.Sha,
) (common type_blobs.Blob, n int64, err error) {
	switch tipe.String() {
	case "", type_blobs.TypeV0:
		store := a.toml_v0
		var blob *type_blobs.V0

		if blob, err = store.GetBlob(blobSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		common = blob

	case type_blobs.TypeV1:
		store := a.toml_v1
		var blob *type_blobs.TomlV1

		if blob, err = store.GetBlob(blobSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		common = blob
	}

	return
}

func (a TypeStore) PutTypedBlob(
	tipe interfaces.ObjectId,
	common type_blobs.Blob,
) (err error) {
	switch tipe.String() {
	case "", type_blobs.TypeV0:
		if blob, ok := common.(*type_blobs.V0); !ok {
			err = errors.Errorf("expected %T but got %T", blob, common)
			return
		} else {
			a.toml_v0.PutBlob(blob)
		}

	case type_blobs.TypeV1:
		if blob, ok := common.(*type_blobs.TomlV1); !ok {
			err = errors.Errorf("expected %T but got %T", blob, common)
			return
		} else {
			a.toml_v1.PutBlob(blob)
		}
	}

	return
}
