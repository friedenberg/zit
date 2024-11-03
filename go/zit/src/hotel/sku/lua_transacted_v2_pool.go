package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool2"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
)

type LuaVMV2 struct {
	lua.LValue
	*lua.VM
	TablePool LuaTablePoolV2
	Selbst    *Transacted
}

func PushTopFuncV2(
	lvm LuaVMPoolV2,
	args []string,
) (vm *LuaVMV2, argsOut []string, err error) {
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
	LuaVMPoolV2    interfaces.Pool2[LuaVMV2, *LuaVMV2]
	LuaTablePoolV2 = interfaces.Pool[LuaTableV2, *LuaTableV2]
)

func MakeLuaVMPoolV2(lvp *lua.VMPool, selbst *Transacted) LuaVMPoolV2 {
	return pool2.MakePool(
		func() (out *LuaVMV2, err error) {
			var vm *lua.VM

			if vm, err = lvp.Pool2.Get(); err != nil {
				err = errors.Wrap(err)
				return
			}

			out = &LuaVMV2{
				VM:        vm,
				TablePool: MakeLuaTablePoolV2(vm),
				Selbst:    selbst,
			}

			return
		},
		nil,
	)
}

func MakeLuaTablePoolV2(vm *lua.VM) LuaTablePoolV2 {
	return pool.MakePool(
		func() (t *LuaTableV2) {
			t = &LuaTableV2{
				Transacted:   vm.Pool.Get(),
				Tags:         vm.Pool.Get(),
				TagsImplicit: vm.Pool.Get(),
			}

			vm.SetField(t.Transacted, "Etiketten", t.Tags)
			vm.SetField(t.Transacted, "EtikettenImplicit", t.TagsImplicit)

			return
		},
		func(t *LuaTableV2) {
			lua.ClearTable(vm.LState, t.Tags)
			lua.ClearTable(vm.LState, t.TagsImplicit)
		},
	)
}
