package sku_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool2"
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
	if vm, err = lvm.Get(); err != nil {
		err = errors.Wrap(err)
		return
	}

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
	LuaVMPool    interfaces.Pool2[LuaVM, *LuaVM]
	LuaTablePool = interfaces.Pool[LuaTable, *LuaTable]
)

func MakeLuaVMPool(lvp *lua.VMPool, selbst *sku.Transacted) LuaVMPool {
	return pool2.MakePool(
		func() (out *LuaVM, err error) {
      var vm *lua.VM

			if vm, err = lvp.Pool2.Get(); err != nil {
        err = errors.Wrap(err)
        return
      }

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
