package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/delta/lua"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/src/juliett/query"
	"code.linenisgreat.com/zit/src/mike/store"
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
	if len(args) == 0 {
		err = errors.Normalf("needs etikett and function name")
		return
	}

	var e kennung.Etikett
	etikett, args := args[0], args[1:]

	if err = e.Set(etikett); err != nil {
		err = errors.Wrap(err)
		return
	}

	funcName, args := args[0], args[1:]

	if len(args) == 0 {
		err = errors.Normalf("function name")
		return
	}

	var lvp store.LuaVMPool

	if lvp, err = u.GetStore().ReadOneSigilLua(e, kennung.SigilCwd); err != nil {
		err = errors.Wrap(err)
		return
	}

	var vm store.LuaVM

	if vm, err = lvp.Get(); err != nil {
		err = errors.Wrap(err)
		return
	}

	fMaybe := vm.GetField(vm.LTable, funcName)

	var f *lua.LFunction
	var ok bool

	if f, ok = fMaybe.(*lua.LFunction); !ok {
		err = errors.Errorf("no such function: %q", funcName)
		return
	}

	vm.Push(f)

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

func (c RunLua) fetchZettels(
	u *umwelt.Umwelt,
	vm *lua.VM,
	t *lua.LTable,
	e kennung.Etikett,
) (err error) {
	qgb := u.MakeQueryBuilderExcludingHidden(kennung.MakeGattung(gattung.Zettel))
	var qg *query.Group

	if qg, err = qgb.BuildQueryGroup(e.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.GetStore().QueryWithCwd(
		qg,
		func(sk *sku.Transacted) (err error) {
			skut := vm.Pool.Get()
			sku_fmt.ToLuaTable(sk, vm.LState, skut)
			t.Append(skut)
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
