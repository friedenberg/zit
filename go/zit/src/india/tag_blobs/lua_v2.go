package tag_blobs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func MakeLuaSelfApplyV2(
	selfOriginal *sku.Transacted,
) interfaces.FuncIter[*lua.VM] {
	if selfOriginal == nil {
		panic("self was nil")
	}

	self := selfOriginal.CloneTransacted()

	return func(vm *lua.VM) (err error) {
		selfTable := sku.MakeLuaTablePoolV2(vm).Get()
		sku.ToLuaTableV2(self, vm.LState, selfTable)
		vm.SetGlobal("Self", selfTable.Transacted)
		return
	}
}

type LuaV2 struct {
	sku.LuaVMPoolV2
}

func (a *LuaV2) GetQueryable() sku.Queryable {
	return a
}

func (a *LuaV2) Reset() {
}

func (a *LuaV2) ResetWith(b LuaV2) {
}

func (tb *LuaV2) ContainsSku(tg sku.TransactedGetter) bool {
	// lb := b.luaVMPoolBuilder.Clone().WithApply(MakeSelfApply(sk))
	vm, err := tb.Get()
	if err != nil {
		ui.Err().Printf("lua script error: %s", err)
		return false
	}

	defer tb.Put(vm)

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

	sku.ToLuaTableV2(
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
