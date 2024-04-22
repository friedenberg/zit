package lua

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/pool"
	lua "github.com/yuin/gopher-lua"
	lua_ast "github.com/yuin/gopher-lua/ast"
	lua_parse "github.com/yuin/gopher-lua/parse"
)

type VM struct {
	*lua.LState
	*lua.LTable
	schnittstellen.Pool[LTable, *LTable]
}

func MakeVMPool(script string) (ml *VMPool, err error) {
	ml = &VMPool{}

	if err = ml.Set(script); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type VMPool struct {
	schnittstellen.Pool[VM, *VM]
}

func (sp *VMPool) Set(script string) (err error) {
	reader := strings.NewReader(script)

	if err = sp.SetReader(reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (sp *VMPool) SetReader(reader io.Reader) (err error) {
	var chunks []lua_ast.Stmt

	if chunks, err = lua_parse.Parse(reader, ""); err != nil {
		err = errors.Wrap(err)
		return
	}

	var compiled *lua.FunctionProto

	if compiled, err = lua.Compile(chunks, ""); err != nil {
		err = errors.Wrap(err)
		return
	}

	sp.Pool = pool.MakePool(
		func() (vm *VM) {
			vm = &VM{
				LState: lua.NewState(),
			}

			vm.Pool = pool.MakePool(
				func() (t *lua.LTable) {
					t = vm.NewTable()
					return
				},
				func(t *lua.LTable) {
					// TODO reset table
				},
			)

			lfunc := vm.NewFunctionFromProto(compiled)
			vm.Push(lfunc)
			errors.PanicIfError(vm.PCall(0, 1, nil))

			retval := vm.LState.Get(1)
			vm.Pop(1)

			if retvalTable, ok := retval.(*LTable); ok {
				vm.LTable = retvalTable
			}

			return vm
		},
		func(vm *VM) {
			vm.SetTop(0)
		},
	)

	return
}
