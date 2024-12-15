package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type ExecLua struct {
	*env.Local
}

func (u ExecLua) Run(sk *sku.Transacted, args ...string) (err error) {
	var lvp sku.LuaVMPoolV1

	if lvp, err = u.GetStore().MakeLuaVMPoolV1WithSku(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	var vm *sku.LuaVMV1

	if vm, args, err = sku.PushTopFuncV1(lvp, args); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, arg := range args {
		vm.Push(lua.LString(arg))
	}

	if err = vm.PCall(len(args), 0, nil); err != nil {
		err = errors.Wrap(err)
		return
	}

	retval := vm.LState.Get(1)

	if retval.Type() != lua.LTNil {
		err = errors.Errorf("lua error: %s", retval)
		return
	}

	return
}
