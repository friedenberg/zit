package query

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func MakeSelfApply(
	self *sku.Transacted,
) interfaces.FuncIter[*lua.VM] {
	if self == nil {
		return nil
	}

	return func(vm *lua.VM) (err error) {
		selbstTable := sku.MakeLuaTablePoolV1(vm).Get()
		sku.ToLuaTableV1(self, vm.LState, selbstTable)
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
	self *sku.Transacted,
) (l Lua, err error) {
	b = b.Clone().WithApply(MakeSelfApply(self))

	var vmp *lua.VMPool

	if vmp, err = b.Build(); err != nil {
		err = errors.Wrap(err)
		return
	}

	l.LuaVMPoolV1 = sku.MakeLuaVMPoolV1(vmp, self)

	return
}

type Lua struct {
	sku.LuaVMPoolV1
}

func (matcher Lua) ContainsSku(tg sku.TransactedGetter) bool {
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

	sku.ToLuaTableV1(
		tg,
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
