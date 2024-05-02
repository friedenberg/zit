package lua

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/pool"
	lua "github.com/yuin/gopher-lua"
)

func MakeVMPool(
	script string,
	require LGFunction,
) (ml *VMPool, err error) {
	ml = &VMPool{}

	if err = ml.Set(script, require); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type VMPool struct {
	schnittstellen.Pool[VM, *VM]
}

func (sp *VMPool) Set(script string, require LGFunction) (err error) {
	reader := strings.NewReader(script)

	if err = sp.SetReader(reader, require); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (sp *VMPool) SetReader(
	reader io.Reader,
	require lua.LGFunction,
) (err error) {
	var compiled *lua.FunctionProto

	if compiled, err = CompileReader(reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	sp.Pool = pool.MakePool(
		func() (vm *VM) {
			vm = &VM{
				LState: lua.NewState(),
			}

			vm.PreloadModule("zit", func(s *lua.LState) int {
				// register functions to the table
				mod := s.SetFuncs(s.NewTable(), map[string]lua.LGFunction{
					"require": require,
				})

				s.Push(mod)

				return 1
			})

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
			errors.PanicIfError(vm.PCall(0, 1, nil))

			vm.LValue = vm.LState.Get(1)
			vm.Pop(1)

			return vm
		},
		func(vm *VM) {
			vm.SetTop(0)
		},
	)

	return
}
