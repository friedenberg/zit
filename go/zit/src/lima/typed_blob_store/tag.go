package typed_blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/tag_blobs"
	"code.linenisgreat.com/zit/go/zit/src/lima/env_lua"
)

type Tag struct {
	dirLayout env_repo.Env
	envLua    env_lua.Env
	toml_v0   TypedStore[tag_blobs.V0, *tag_blobs.V0]
	toml_v1   TypedStore[tag_blobs.TomlV1, *tag_blobs.TomlV1]
	lua_v1    TypedStore[tag_blobs.LuaV1, *tag_blobs.LuaV1]
	lua_v2    TypedStore[tag_blobs.LuaV2, *tag_blobs.LuaV2]
}

func MakeTagStore(
	dirLayout env_repo.Env,
	envLua env_lua.Env,
) Tag {
	return Tag{
		dirLayout: dirLayout,
		envLua:    envLua,
		toml_v0: MakeBlobStore(
			dirLayout,
			MakeBlobFormat(
				MakeTomlDecoderIgnoreTomlErrors[tag_blobs.V0](
					dirLayout,
				),
				TomlBlobEncoder[tag_blobs.V0, *tag_blobs.V0]{},
				dirLayout,
			),
			func(a *tag_blobs.V0) {
				a.Reset()
			},
		),
		toml_v1: MakeBlobStore(
			dirLayout,
			MakeBlobFormat(
				MakeTomlDecoderIgnoreTomlErrors[tag_blobs.TomlV1](
					dirLayout,
				),
				TomlBlobEncoder[tag_blobs.TomlV1, *tag_blobs.TomlV1]{},
				dirLayout,
			),
			func(a *tag_blobs.TomlV1) {
				a.Reset()
			},
		),
		lua_v1: MakeBlobStore(
			dirLayout,
			MakeBlobFormat[tag_blobs.LuaV1, *tag_blobs.LuaV1](
				nil,
				nil,
				dirLayout,
			),
			func(a *tag_blobs.LuaV1) {
			},
		),
		lua_v2: MakeBlobStore(
			dirLayout,
			MakeBlobFormat[tag_blobs.LuaV2, *tag_blobs.LuaV2](
				nil,
				nil,
				dirLayout,
			),
			func(a *tag_blobs.LuaV2) {
			},
		),
	}
}

func (a Tag) GetCommonStore() sku.BlobStore[tag_blobs.Blob] {
	return a
}

func (a Tag) GetTransactedWithBlob(
	tg sku.TransactedGetter,
) (twb sku.TransactedWithBlob[tag_blobs.Blob], n int64, err error) {
	sk := tg.GetSku()
	tipe := sk.GetType()
	blobSha := sk.GetBlobSha()

	twb.Transacted = sk.CloneTransacted()

	switch tipe.String() {
	case "", builtin_types.TagTypeTomlV0:
		store := a.toml_v0
		var blob *tag_blobs.V0

		if blob, err = store.GetBlob(blobSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		twb.Blob = blob

	case builtin_types.TagTypeTomlV1:
		store := a.toml_v1
		var blob *tag_blobs.TomlV1

		if blob, err = store.GetBlob(blobSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		lb := a.envLua.MakeLuaVMPoolBuilder().WithApply(tag_blobs.MakeLuaSelfApplyV1(sk))

		var vmp *lua.VMPool

		lb.WithScript(blob.Filter)

		if vmp, err = lb.Build(); err != nil {
			err = errors.Wrap(err)
			return
		}

		blob.LuaVMPoolV1 = sku.MakeLuaVMPoolV1(vmp, nil)
		twb.Blob = blob

	case builtin_types.TagTypeLuaV1:
		var rc sha.ReadCloser

		if rc, err = a.dirLayout.BlobReader(blobSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, rc)

		lb := a.envLua.MakeLuaVMPoolBuilder().WithApply(tag_blobs.MakeLuaSelfApplyV1(sk))

		var vmp *lua.VMPool

		lb.WithReader(rc)

		if vmp, err = lb.Build(); err != nil {
			err = errors.Wrap(err)
			return
		}

		twb.Blob = &tag_blobs.LuaV1{
			LuaVMPoolV1: sku.MakeLuaVMPoolV1(vmp, nil),
		}

	case builtin_types.TagTypeLuaV2:
		var rc sha.ReadCloser

		if rc, err = a.dirLayout.BlobReader(blobSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, rc)

		lb := a.envLua.MakeLuaVMPoolBuilder().WithApply(tag_blobs.MakeLuaSelfApplyV2(sk))

		var vmp *lua.VMPool

		lb.WithReader(rc)

		if vmp, err = lb.Build(); err != nil {
			err = errors.Wrap(err)
			return
		}

		twb.Blob = &tag_blobs.LuaV2{
			LuaVMPoolV2: sku.MakeLuaVMPoolV2(vmp, nil),
		}
	}

	return
}

func (a Tag) PutTransactedWithBlob(
	twb sku.TransactedWithBlob[tag_blobs.Blob],
) (err error) {
	tipe := twb.GetType()

	switch tipe.String() {
	case "", builtin_types.TagTypeTomlV0:
		if blob, ok := twb.Blob.(*tag_blobs.V0); !ok {
			err = errors.ErrorWithStackf("expected %T but got %T", blob, twb.Blob)
			return
		} else {
			a.toml_v0.PutBlob(blob)
		}

	case builtin_types.TagTypeLuaV1:
		if blob, ok := twb.Blob.(*tag_blobs.TomlV1); !ok {
			err = errors.ErrorWithStackf("expected %T but got %T", blob, twb.Blob)
			return
		} else {
			a.toml_v1.PutBlob(blob)
		}
	}

	sku.GetTransactedPool().Put(twb.Transacted)

	return
}
