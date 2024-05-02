package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/lua"
	"code.linenisgreat.com/zit/src/november/umwelt"
)

type RunLua struct{}

func init() {
	registerCommand(
		"run-lua",
		func(f *flag.FlagSet) Command {
			c := &RunLua{}

			return c
		},
	)
}

func (c RunLua) Run(u *umwelt.Umwelt, args ...string) (err error) {
	script := args[0]
	funcName := args[1]

	var vp *lua.VMPool

	if vp, err = u.GetStore().MakeLuaVMPool(script); err != nil {
		err = errors.Wrap(err)
		return
	}

	vm := vp.Get()
	defer vp.Put(vm)

	var t *lua.LTable

	if t, err = vm.GetTopTableOrError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	f := vm.GetField(t, funcName).(*lua.LFunction)
	vm.Push(f)

	if err = vm.PCall(0, 0, nil); err != nil {
		err = errors.Wrap(err)
		return
	}

	retval := vm.LState.Get(1)
	// vm.Pop(1)

	if retval.Type() != lua.LTNil {
		err = errors.Errorf("lua error: %s", retval)
		return
	}

	return
}
