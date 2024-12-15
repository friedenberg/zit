package env

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (u *Local) MakeLuaVMPoolBuilder() *lua.VMPoolBuilder {
	return (&lua.VMPoolBuilder{}).WithSearcher(u.LuaSearcher)
}

func (s *Local) GetSkuFromString(lv string) (sk *sku.Transacted, err error) {
	sk = sku.GetTransactedPool().Get()

	defer func() {
		if err != nil {
			return
		}

		if err = s.GetStore().ReadOneInto(sk.GetObjectId(), sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}()

	if err = sk.ObjectId.SetOnlyNotUnknownGenre(lv); err == nil {
		return
	}

	rb := catgut.MakeRingBuffer(strings.NewReader(lv), 0)

	if _, err = s.luaSkuFormat.ReadStringFormat(
		catgut.MakeRingBufferRuneScanner(rb),
		sk,
	); err == nil {
		return
	}

	return
}

func (s *Local) LuaSearcher(ls *lua.LState) int {
	lv := ls.ToString(1)
	ls.Pop(1)

	var err error
	var sk *sku.Transacted

	if sk, err = s.GetSkuFromString(lv); err != nil {
		ls.Push(lua.LString(err.Error()))
		return 1
	}

	sku.GetTransactedPool().Put(sk)

	ls.Push(ls.NewFunction(s.LuaRequire))

	return 1
}

// TODO modify `package.loaded` to include variations of object id
func (s *Local) LuaRequire(ls *lua.LState) int {
	// TODO handle second extra arg
	// TODO parse lv as object id / blob
	lv := ls.ToString(1)
	ls.Pop(1)

	var err error
	var sk *sku.Transacted

	if sk, err = s.GetSkuFromString(lv); err != nil {
		panic(err)
		// ls.Push(lua.LString(err.Error()))
		// return 1
	}

	defer sku.GetTransactedPool().Put(sk)

	if err = s.GetStore().ReadOneInto(sk.GetObjectId(), sk); err != nil {
		panic(err)
	}

	var ar sha.ReadCloser

	if ar, err = s.GetStore().GetDirectoryLayout().BlobReader(
		sk.GetBlobSha(),
	); err != nil {
		panic(err)
	}

	defer errors.DeferredCloser(&err, ar)

	var compiled *lua.FunctionProto

	if compiled, err = lua.CompileReader(ar); err != nil {
		panic(err)
	}

	ls.Push(ls.NewFunctionFromProto(compiled))

	if err = ls.PCall(0, 1, nil); err != nil {
		panic(err)
	}

	return 1
}
