package lua

import (
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/pool"
	lua "github.com/yuin/gopher-lua"
)

type VM struct {
	*lua.LState
	lua.LValue
	schnittstellen.Pool[LTable, *LTable]
}

func MakeVMFromScript(
	require LGFunction,
	script string,
) (vm *VM, err error) {
	reader := strings.NewReader(script)

	var compiled *lua.FunctionProto

	if compiled, err = CompileReader(reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	if vm, err = MakeVMFromCompiled(require, compiled); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeVMFromCompiled(
	require LGFunction,
	compiled *lua.FunctionProto,
) (vm *VM, err error) {
	vm = &VM{
		LState: lua.NewState(),
	}

	// vm.PreloadModule("zit", func(s *lua.LState) int {
	// 	// register functions to the table
	// 	mod := s.SetFuncs(s.NewTable(), map[string]lua.LGFunction{
	// 		"require": require,
	// 	})

	// 	s.Push(mod)

	// 	return 1
	// })

	vm.Pool = pool.MakePool(
		func() (t *lua.LTable) {
			t = vm.NewTable()
			return
		},
		func(t *lua.LTable) {
			// TODO reset table
		},
	)

	tableZit := vm.Pool.Get()
	vm.SetField(tableZit, "require", vm.NewFunction(require))
	vm.SetGlobal("zit", tableZit)

	lfunc := vm.NewFunctionFromProto(compiled)
	vm.Push(lfunc)

	if err = vm.PCall(0, 1, nil); err != nil {
		err = errors.Wrap(err)
		return
	}

	vm.LValue = vm.LState.Get(1)
	vm.Pop(1)

	return
}

func (vm *VM) GetTopTableOrError() (t *LTable, err error) {
	if vm.LValue.Type() != LTTable {
		err = errors.Errorf("expected %v but got %v", LTTable, vm.LValue.Type())
		return
	}

	t = vm.LValue.(*LTable)

	return
}

func (vm *VM) GetTopFunctionOrError() (t *LFunction, err error) {
	if vm.LValue.Type() != LTFunction {
		err = errors.Errorf("expected %v but got %v", LTFunction, vm.LValue.Type())
		return
	}

	t = vm.LValue.(*LFunction)

	return
}
