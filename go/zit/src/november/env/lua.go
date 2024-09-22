package env

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Env) GetSkuFromString(lv string) (sk *sku.Transacted, err error) {
	e := sku.GetTransactedPool().Get()
	sk = e.GetSku()

	defer func() {
		if err != nil {
			return
		}

		if err = s.GetStore().ReadOneInto(sk.GetObjectId(), sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}()

	if err = sk.ObjectId.Set(lv); err == nil {
		return
	}

	rb := catgut.MakeRingBuffer(strings.NewReader(lv), 0)

	if _, err = s.luaSkuFormat.ReadStringFormat(rb, e); err == nil {
		return
	}

	return
}

func (s *Env) LuaSearcher(ls *lua.LState) int {
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
func (s *Env) LuaRequire(ls *lua.LState) int {
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

	if ar, err = s.GetStore().GetStandort().BlobReader(
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

// TODO check if selbst needs to be passed in
func (u *Env) MakeLuaVMPool(script string) (vp *lua.VMPool, err error) {
	if vp, err = lua.MakeVMPoolWithZitSearcher(
		script,
		u.LuaSearcher,
		nil,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
