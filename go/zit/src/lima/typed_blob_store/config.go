package typed_blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_mutable_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

type Config struct {
	toml_v0 TypedStore[config_mutable_blobs.V0, *config_mutable_blobs.V0]
	toml_v1 TypedStore[config_mutable_blobs.V1, *config_mutable_blobs.V1]
}

func MakeConfigStore(
	repoLayout env_repo.Env,
) Config {
	return Config{
		toml_v0: MakeBlobStore(
			repoLayout,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[config_mutable_blobs.V0](
					repoLayout,
				),
				ParsedBlobTomlFormatter[config_mutable_blobs.V0, *config_mutable_blobs.V0]{},
				repoLayout,
			),
			func(a *config_mutable_blobs.V0) {
				a.Reset()
			},
		),
		toml_v1: MakeBlobStore(
			repoLayout,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[config_mutable_blobs.V1](
					repoLayout,
				),
				ParsedBlobTomlFormatter[config_mutable_blobs.V1, *config_mutable_blobs.V1]{},
				repoLayout,
			),
			func(a *config_mutable_blobs.V1) {
				a.Reset()
			},
		),
	}
}

func (a Config) ParseTypedBlob(
	tipe ids.Type,
	blobSha interfaces.Sha,
) (common config_mutable_blobs.Blob, n int64, err error) {
	switch tipe.String() {
	case "", builtin_types.ConfigTypeTomlV0:
		store := a.toml_v0
		var blob *config_mutable_blobs.V0

		if blob, err = store.GetBlob(blobSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		common = blob

	case builtin_types.ConfigTypeTomlV1:
		store := a.toml_v1
		var blob *config_mutable_blobs.V1

		if blob, err = store.GetBlob(blobSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		common = blob
	}

	return
}

func (a Config) FormatTypedBlob(
	tg sku.TransactedGetter,
	w io.Writer,
) (n int64, err error) {
	sk := tg.GetSku()

	tipe := sk.GetType()
	blobSha := sk.GetBlobSha()

	var store interfaces.SavedBlobFormatter
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

func (a Config) PutTypedBlob(
	tipe interfaces.ObjectId,
	common config_mutable_blobs.Blob,
) (err error) {
	switch tipe.String() {
	case "", builtin_types.ConfigTypeTomlV0:
		if blob, ok := common.(*config_mutable_blobs.V0); !ok {
			err = errors.Errorf("expected %T but got %T", blob, common)
			return
		} else {
			a.toml_v0.PutBlob(blob)
		}

	case builtin_types.ConfigTypeTomlV1:
		if blob, ok := common.(*config_mutable_blobs.V1); !ok {
			err = errors.Errorf("expected %T but got %T", blob, common)
			return
		} else {
			a.toml_v1.PutBlob(blob)
		}
	}

	return
}
