package lua

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/pool"
	lua "github.com/yuin/gopher-lua"
)

func MakeVMPoolWithZitRequire(
	script string,
	require LGFunction,
) (ml *VMPool, err error) {
	ml = &VMPool{
		Require: require,
	}

	if err = ml.Set(script); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeVMPoolWithZitSearcher(
	script string,
	searcher LGFunction,
) (ml *VMPool, err error) {
	ml = &VMPool{
		Searcher: searcher,
	}

	if err = ml.Set(script); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type VMPool struct {
	schnittstellen.Pool[VM, *VM]
	Require  LGFunction
	Searcher LGFunction
	compiled *lua.FunctionProto
}

func (sp *VMPool) Set(script string) (err error) {
	reader := strings.NewReader(script)

	if err = sp.SetReader(reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (sp *VMPool) PrepareVM(vm *VM) {
	vm.Pool = pool.MakePool(
		func() (t *lua.LTable) {
			t = vm.NewTable()
			return
		},
		func(t *lua.LTable) {
			// TODO reset table
		},
	)

	if sp.Require != nil {
		vm.PreloadModule("zit", func(s *lua.LState) int {
			// register functions to the table
			mod := s.SetFuncs(s.NewTable(), map[string]lua.LGFunction{
				"require": sp.Require,
			})

			s.Push(mod)

			return 1
		})

		tableZit := vm.Pool.Get()
		vm.SetField(tableZit, "require", vm.NewFunction(sp.Require))
		vm.SetGlobal("zit", tableZit)
	}

	if sp.Searcher != nil {
		packageTable := vm.GetGlobal("package").(*LTable)

		if true { // lua <= 5.1
			loaderTable := vm.GetField(packageTable, "loaders").(*LTable)
			loaderTable.Insert(1, vm.NewFunction(sp.Searcher))
		} else {
			searcherTable := vm.Pool.Get()
			packageTable.Insert(1, searcherTable)
			searcherTable.Insert(1, vm.NewFunction(sp.Searcher))
		}
	}

	lfunc := vm.NewFunctionFromProto(sp.compiled)
	vm.Push(lfunc)
	errors.PanicIfError(vm.PCall(0, 1, nil))

	vm.LValue = vm.LState.Get(1)
	vm.Pop(1)
}

func (sp *VMPool) SetReader(
	reader io.Reader,
) (err error) {
	if sp.compiled, err = CompileReader(reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	sp.Pool = pool.MakePool(
		func() (vm *VM) {
			vm = &VM{
				LState: lua.NewState(),
			}

			sp.PrepareVM(vm)

			return vm
		},
		func(vm *VM) {
			vm.SetTop(0)
		},
	)

	return
}
