package typed_blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/type_blobs"
)

type Type struct {
	toml_v0 TypedStore[type_blobs.V0, *type_blobs.V0]
	toml_v1 TypedStore[type_blobs.TomlV1, *type_blobs.TomlV1]
}

func MakeTypeStore(
	envRepo env_repo.Env,
) Type {
	return Type{
		toml_v0: MakeBlobStore(
			envRepo,
			MakeBlobFormat(
				MakeTomlDecoderIgnoreTomlErrors[type_blobs.V0](
					envRepo,
				),
				TomlBlobEncoder[type_blobs.V0, *type_blobs.V0]{},
				envRepo,
			),
			func(a *type_blobs.V0) {
				a.Reset()
			},
		),
		toml_v1: MakeBlobStore(
			envRepo,
			MakeBlobFormat(
				MakeTomlDecoderIgnoreTomlErrors[type_blobs.TomlV1](
					envRepo,
				),
				TomlBlobEncoder[type_blobs.TomlV1, *type_blobs.TomlV1]{},
				envRepo,
			),
			func(a *type_blobs.TomlV1) {
				a.Reset()
			},
		),
	}
}

func (a Type) GetCommonStore() interfaces.TypedBlobStore[type_blobs.Blob] {
	return a
}

func (a Type) ParseTypedBlob(
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

func (a Type) PutTypedBlob(
	tipe interfaces.ObjectId,
	common type_blobs.Blob,
) (err error) {
	switch tipe.String() {
	case "", type_blobs.TypeV0:
		if blob, ok := common.(*type_blobs.V0); !ok {
			err = errors.ErrorWithStackf("expected %T but got %T", blob, common)
			return
		} else {
			a.toml_v0.PutBlob(blob)
		}

	case type_blobs.TypeV1:
		if blob, ok := common.(*type_blobs.TomlV1); !ok {
			err = errors.ErrorWithStackf("expected %T but got %T", blob, common)
			return
		} else {
			a.toml_v1.PutBlob(blob)
		}
	}

	return
}
