package lua

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	lua "github.com/yuin/gopher-lua"
)

type VM struct {
	*lua.LState
	Top lua.LValue
	interfaces.Pool[LTable, *LTable]
}

func (vm *VM) GetTopFunctionOrFunctionNamedError(
	args []string,
) (t *LFunction, argsOut []string, err error) {
	funcName := ""
	argsOut = args

	if len(args) > 0 {
		funcName = args[0]
	}

	if vm.Top.Type() == LTTable {
		if funcName == "" {
			err = errors.Errorf("needs function name because top is table")
			return
		}

		tt := vm.Top.(*LTable)

		tv := vm.GetField(tt, funcName)

		if tv.Type() != LTFunction {
			err = errors.Errorf("expected %v but got %v", LTFunction, tv.Type())
			return
		}

		argsOut = argsOut[1:]

		t = tv.(*LFunction)
	} else if vm.Top.Type() == LTFunction {
		t = vm.Top.(*LFunction)
	} else {
		err = errors.Errorf(
			"expected table or function but got: %q",
			vm.Top.Type(),
		)

		return
	}

	return
}

func (vm *VM) GetTopTableOrError() (t *LTable, err error) {
	if vm.Top.Type() != LTTable {
		err = errors.Errorf("expected %v but got %v", LTTable, vm.Top.Type())
		return
	}

	t = vm.Top.(*LTable)

	return
}

func (vm *VM) GetTopFunctionOrError() (t *LFunction, err error) {
	if vm.Top.Type() != LTFunction {
		err = errors.Errorf("expected %v but got %v", LTFunction, vm.Top.Type())
		return
	}

	t = vm.Top.(*LFunction)

	return
}
