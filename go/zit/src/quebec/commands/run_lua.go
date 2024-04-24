package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/lua"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/src/juliett/query"
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
	args = args[1:]

	b := u.MakeQueryBuilderExcludingHidden(kennung.MakeGattung())

	var qg *query.Group

	if qg, err = b.BuildQueryGroup(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	var vp lua.VMPool

	if err = vp.Set(script); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := u.PrinterTransactedLike()

	if err = u.GetStore().QueryWithCwd(
		qg,
		func(sk *sku.Transacted) (err error) {
			vm := vp.Get()
			defer vp.Put(vm)

			var f *lua.LFunction

			if f, err = vm.GetTopFunctionOrError(); err != nil {
				err = errors.Wrap(err)
				return
			}

			tableKinder := vm.Pool.Get()
			defer vm.Put(tableKinder)

			sku_fmt.ToLuaTable(
				sk,
				vm.LState,
				tableKinder,
			)

			vm.Push(f)
			vm.Push(tableKinder)
			vm.Call(
				1,
				1,
			)

			retval := vm.LState.Get(1)
			vm.Pop(1)

			if retval.Type() != lua.LTNil {
				err = errors.Errorf("lua error: %s", retval)
				return
			}

			if err = sku_fmt.FromLuaTable(sk, vm.LState, tableKinder); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = p(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
