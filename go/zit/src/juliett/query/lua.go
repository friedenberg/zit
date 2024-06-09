package query

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/src/delta/lua"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
)

func MakeSelbstApply(
	selbst *sku.Transacted,
) schnittstellen.FuncIter[*lua.VM] {
	if selbst == nil {
		return nil
	}

	return func(vm *lua.VM) (err error) {
		selbstTable := vm.Pool.Get()
		sku_fmt.ToLuaTable(selbst, vm.LState, selbstTable)
		vm.SetGlobal("Selbst", selbstTable)
		return
	}
}

func MakeLua(
	selbst *sku.Transacted,
	script string,
	require lua.LGFunction,
) (ml Lua, err error) {
	ml = Lua{
		VMPool: (&lua.VMPoolBuilder{}).WithRequire(require).Build(),
	}

	if err = ml.Set(script, MakeSelbstApply(selbst)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type Lua struct {
	*lua.VMPool
}

func (matcher Lua) ContainsSku(sk *sku.Transacted) bool {
	vm := matcher.Get()
	defer matcher.Put(vm)

	var t *lua.LTable
	var err error

	t, err = vm.GetTopTableOrError()
	if err != nil {
		ui.Err().Print(err)
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
		ui.Err().Print(err)
		return false
	}

	retval := vm.LState.Get(1)
	vm.Pop(1)

	if retval.Type() != lua.LTBool {
		ui.Err().Printf("expected bool but got %s", retval.Type())
		return false
	}

	return bool(retval.(lua.LBool))
}
