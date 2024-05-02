package query

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/delta/lua"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
)

func MakeLua(script string, require lua.LGFunction) (ml *Lua, err error) {
	ml = &Lua{}

	if err = ml.Set(script, require); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type Lua struct {
	lua.VMPool
}

func (matcher *Lua) ContainsSku(sk *sku.Transacted) bool {
	vm := matcher.Get()
	defer matcher.Put(vm)

	var t *lua.LTable
	var err error

	t, err = vm.GetTopTableOrError()
	if err != nil {
		log.Err().Print(err)
		return false
	}

	f := vm.GetField(t, "contains_sku").(*lua.LFunction)

	tSku := vm.Pool.Get()
	defer vm.Put(tSku)

	vm.Push(f)

	sku_fmt.ToLuaTable(
		sk,
		vm.LState,
		tSku,
	)

	vm.Push(tSku)

	err = vm.PCall(1, 1, nil)
	if err != nil {
		log.Err().Print(err)
		return false
	}

	retval := vm.LState.Get(1)
	vm.Pop(1)

	if retval.Type() != lua.LTBool {
		log.Err().Printf("expected bool but got %s", retval.Type())
		return false
	}

	return bool(retval.(lua.LBool))
}
