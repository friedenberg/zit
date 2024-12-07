package blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/tag_blobs"
)

type TagStore struct {
	dirLayout        dir_layout.DirLayout
	luaVMPoolBuilder *lua.VMPoolBuilder
	toml_v0          Store[tag_blobs.V0, *tag_blobs.V0]
	toml_v1          Store[tag_blobs.TomlV1, *tag_blobs.TomlV1]
	lua_v1           Store[tag_blobs.LuaV1, *tag_blobs.LuaV1]
	lua_v2           Store[tag_blobs.LuaV2, *tag_blobs.LuaV2]
}

func MakeTagStore(
	dirLayout dir_layout.DirLayout,
	luaVMPoolBuilder *lua.VMPoolBuilder,
) TagStore {
	return TagStore{
		dirLayout:        dirLayout,
		luaVMPoolBuilder: luaVMPoolBuilder,
		toml_v0: MakeBlobStore(
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
		toml_v1: MakeBlobStore(
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

func (a TagStore) GetCommonStore() CommonStore2[tag_blobs.Blob] {
	return a
}

func (a TagStore) GetTransactedWithBlob(
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

		lb := a.luaVMPoolBuilder.Clone().WithApply(tag_blobs.MakeLuaSelfApplyV1(sk))

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

		lb := a.luaVMPoolBuilder.Clone().WithApply(tag_blobs.MakeLuaSelfApplyV1(sk))

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

		lb := a.luaVMPoolBuilder.Clone().WithApply(tag_blobs.MakeLuaSelfApplyV2(sk))

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

func (a TagStore) PutTransactedWithBlob(
	twb sku.TransactedWithBlob[tag_blobs.Blob],
) (err error) {
	tipe := twb.GetType()

	switch tipe.String() {
	case "", builtin_types.TagTypeTomlV0:
		if blob, ok := twb.Blob.(*tag_blobs.V0); !ok {
			err = errors.Errorf("expected %T but got %T", blob, twb.Blob)
			return
		} else {
			a.toml_v0.PutBlob(blob)
		}

	case builtin_types.TagTypeLuaV1:
		if blob, ok := twb.Blob.(*tag_blobs.TomlV1); !ok {
			err = errors.Errorf("expected %T but got %T", blob, twb.Blob)
			return
		} else {
			a.toml_v1.PutBlob(blob)
		}
	}

	sku.GetTransactedPool().Put(twb.Transacted)

	return
}
