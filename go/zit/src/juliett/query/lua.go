package query

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

func MakeSelfApply(
	self *sku.Transacted,
) interfaces.FuncIter[*lua.VM] {
	if self == nil {
		return nil
	}

	return func(vm *lua.VM) (err error) {
		selbstTable := sku_fmt.MakeLuaTablePool(vm).Get()
		sku_fmt.ToLuaTable(self, vm.LState, selbstTable)
		vm.SetGlobal("Selbst", selbstTable.Transacted)
		return
	}
}

func MakeLua(
	self *sku.Transacted,
	script string,
	require lua.LGFunction,
) (ml Lua, err error) {
	b := (&lua.VMPoolBuilder{}).
		WithScript(script).
		WithRequire(require)

	if ml, err = MakeLuaFromBuilder(b, self); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeLuaFromBuilder(
	b *lua.VMPoolBuilder,
	selbst *sku.Transacted,
) (l Lua, err error) {
	b = b.Clone().WithApply(MakeSelfApply(selbst))

	var vmp *lua.VMPool

	if vmp, err = b.Build(); err != nil {
		err = errors.Wrap(err)
		return
	}

	l.LuaVMPool = sku_fmt.MakeLuaVMPool(vmp, selbst)

	return
}

type Lua struct {
	sku_fmt.LuaVMPool
}

func (matcher Lua) ContainsSku(sk *sku.Transacted) bool {
	vm, err := matcher.Get()

	if err != nil {
		ui.Err().Printf("lua script error: %s", err)
		return false
	}

	defer matcher.Put(vm)

	var t *lua.LTable

	t, err = vm.VM.GetTopTableOrError()
	if err != nil {
		ui.Err().Print(err)
		return false
	}

	// TODO safer
	f := vm.VM.GetField(t, "contains_sku").(*lua.LFunction)

	tSku := vm.TablePool.Get()
	defer vm.TablePool.Put(tSku)

	vm.VM.Push(f)

	sku_fmt.ToLuaTable(
		sk,
		vm.VM.LState,
		tSku,
	)

	vm.VM.Push(tSku.Transacted)

	err = vm.VM.PCall(1, 1, nil)
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
