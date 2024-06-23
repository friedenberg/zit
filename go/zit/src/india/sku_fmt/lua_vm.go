package sku_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type LuaVM struct {
	lua.LValue
	*lua.VM
	TablePool LuaTablePool
	Selbst    *sku.Transacted
}

func PushTopFunc(
	lvm LuaVMPool,
	args []string,
) (vm *LuaVM, argsOut []string, err error) {
	vm = lvm.Get()
	vm.LValue = vm.Top

	var f *lua.LFunction

	if f, argsOut, err = vm.GetTopFunctionOrFunctionNamedError(
		args,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	vm.Push(f)

	return
}

type (
	LuaVMPool    schnittstellen.Pool[LuaVM, *LuaVM]
	LuaTablePool = schnittstellen.Pool[LuaTable, *LuaTable]
)

func MakeLuaVMPool(lvp *lua.VMPool, selbst *sku.Transacted) LuaVMPool {
	return pool.MakePool(
		func() (out *LuaVM) {
			vm := lvp.Pool.Get()
			out = &LuaVM{
				VM:        vm,
				TablePool: MakeLuaTablePool(vm),
				Selbst:    selbst,
			}

			return
		},
		nil,
	)
}

func MakeLuaTablePool(vm *lua.VM) LuaTablePool {
	return pool.MakePool(
		func() (t *LuaTable) {
			t = &LuaTable{
				Transacted:        vm.Pool.Get(),
				Etiketten:         vm.Pool.Get(),
				EtikettenImplicit: vm.Pool.Get(),
			}

			vm.SetField(t.Transacted, "Etiketten", t.Etiketten)
			vm.SetField(t.Transacted, "EtikettenImplicit", t.EtikettenImplicit)

			return
		},
		func(t *LuaTable) {
			lua.ClearTable(vm.LState, t.Etiketten)
			lua.ClearTable(vm.LState, t.EtikettenImplicit)
		},
	)
}
