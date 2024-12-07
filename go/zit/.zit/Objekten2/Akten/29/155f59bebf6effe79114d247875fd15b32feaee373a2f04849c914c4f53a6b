package blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/golf/mutable_config_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type ConfigStore struct {
	toml_v0 Store[mutable_config_blobs.V0, *mutable_config_blobs.V0]
	toml_v1 Store[mutable_config_blobs.V1, *mutable_config_blobs.V1]
}

func MakeConfigStore(
	dirLayout dir_layout.DirLayout,
) ConfigStore {
	return ConfigStore{
		toml_v0: MakeBlobStore(
			dirLayout,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[mutable_config_blobs.V0](
					dirLayout,
				),
				ParsedBlobTomlFormatter[mutable_config_blobs.V0, *mutable_config_blobs.V0]{},
				dirLayout,
			),
			func(a *mutable_config_blobs.V0) {
				a.Reset()
			},
		),
		toml_v1: MakeBlobStore(
			dirLayout,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[mutable_config_blobs.V1](
					dirLayout,
				),
				ParsedBlobTomlFormatter[mutable_config_blobs.V1, *mutable_config_blobs.V1]{},
				dirLayout,
			),
			func(a *mutable_config_blobs.V1) {
				a.Reset()
			},
		),
	}
}

func (a ConfigStore) ParseTypedBlob(
	tipe interfaces.ObjectId,
	blobSha interfaces.Sha,
) (common mutable_config_blobs.Blob, n int64, err error) {
	switch tipe.String() {
	case "", builtin_types.ConfigTypeTomlV0:
		store := a.toml_v0
		var blob *mutable_config_blobs.V0

		if blob, err = store.GetBlob(blobSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		common = blob

	case builtin_types.ConfigTypeTomlV1:
		store := a.toml_v1
		var blob *mutable_config_blobs.V1

		if blob, err = store.GetBlob(blobSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		common = blob
	}

	return
}

func (a ConfigStore) FormatTypedBlob(
	tg sku.TransactedGetter,
	w io.Writer,
) (n int64, err error) {
	sk := tg.GetSku()

	tipe := sk.GetType()
	blobSha := sk.GetBlobSha()

	var store SavedBlobFormatter
	switch tipe.String() {
	case "", builtin_types.ConfigTypeTomlV0:
		store = a.toml_v0

	case builtin_types.ConfigTypeTomlV1:
		store = a.toml_v1
	}

	if n, err = store.FormatSavedBlob(w, blobSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a ConfigStore) PutTypedBlob(
	tipe interfaces.ObjectId,
	common mutable_config_blobs.Blob,
) (err error) {
	switch tipe.String() {
	case "", builtin_types.ConfigTypeTomlV0:
		if blob, ok := common.(*mutable_config_blobs.V0); !ok {
			err = errors.Errorf("expected %T but got %T", blob, common)
			return
		} else {
			a.toml_v0.PutBlob(blob)
		}

	case builtin_types.ConfigTypeTomlV1:
		if blob, ok := common.(*mutable_config_blobs.V1); !ok {
			err = errors.Errorf("expected %T but got %T", blob, common)
			return
		} else {
			a.toml_v1.PutBlob(blob)
		}
	}

	return
}
