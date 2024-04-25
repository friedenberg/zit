package store

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/lua"
	"code.linenisgreat.com/zit/src/delta/sha"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

func (u *Store) LuaRequire(s *lua.LState) int {
	// TODO parse lv as kennung/akte
	lv := s.ToString(1)
	s.Pop(1)

	// // TODO load akte as module if possible
	// s.Push(lua.LString(lv))

	var err error
	var k kennung.Kennung2

	if err = k.Set(lv); err != nil {
		panic(err)
	}

	var sk *sku.Transacted

	// TODO add support for cwd in readonekennung
	if sk, err = u.ReadOneKennung(&k); err != nil {
		panic(err)
	}

	var ar sha.ReadCloser

	if ar, err = u.GetStandort().AkteReader(
		sk.GetAkteSha(),
	); err != nil {
		panic(err)
	}

	defer errors.DeferredCloser(&err, ar)

	var compiled *lua.FunctionProto

	if compiled, err = lua.CompileReader(ar); err != nil {
		panic(err)
	}

	s.Push(s.NewFunctionFromProto(compiled))

	if err = s.PCall(0, 1, nil); err != nil {
		panic(err)
	}

	return 1
}

func (u *Store) MakeLuaVMPool(script string) (vp *lua.VMPool, err error) {
	if vp, err = lua.MakeVMPool(script, u.LuaRequire); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
