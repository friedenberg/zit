package store

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/delta/lua"
	"code.linenisgreat.com/zit/src/delta/typ_akte"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
)

func (s *Store) tryNewHook(
	kinder *sku.Transacted,
	mode objekte_mode.Mode,
) (err error) {
	var t *sku.Transacted

	if t, err = s.ReadOneKennung(kinder.GetTyp()); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	var akte *typ_akte.V0

	if akte, err = s.GetAkten().GetTypV0().GetAkte(t.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	script, ok := akte.Hooks.(string)

	if !ok || script == "" {
		return
	}

	var vp lua.VMPool

	if err = vp.Set(script); err != nil {
		err = errors.Wrap(err)
		return
	}

	vm := vp.Get()
	defer vp.Put(vm)

	f := vm.GetField(vm.LTable, "on_new")

	if f.Type() != lua.LTFunction {
		return
	}

	tableKinder := vm.Pool.Get()
	defer vm.Put(tableKinder)

	sku_fmt.ToLuaTable(
		kinder,
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

	if err = sku_fmt.FromLuaTable(kinder, vm.LState, tableKinder); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) tryPreCommitHook(
	kinder *sku.Transacted,
	mutter *sku.Transacted,
	mode objekte_mode.Mode,
) (err error) {
	var t *sku.Transacted

	if t, err = s.ReadOneKennung(kinder.GetTyp()); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	var akte *typ_akte.V0

	if akte, err = s.GetAkten().GetTypV0().GetAkte(t.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	script, ok := akte.Hooks.(string)

	if !ok || script == "" {
		return
	}

	var vp lua.VMPool

	if err = vp.Set(script); err != nil {
		err = errors.Wrap(err)
		return
	}

	vm := vp.Get()
	defer vp.Put(vm)

	f := vm.GetField(vm.LTable, "on_pre_commit")

	if f.Type() != lua.LTFunction {
		return
	}

	tableKinder := vm.Pool.Get()
	defer vm.Put(tableKinder)

	sku_fmt.ToLuaTable(
		kinder,
		vm.LState,
		tableKinder,
	)

	var tableMutter *lua.LTable

	if mutter != nil {
		tableMutter = vm.Pool.Get()
		defer vm.Put(tableMutter)

		sku_fmt.ToLuaTable(
			mutter,
			vm.LState,
			tableKinder,
		)
	}

	vm.Push(f)
	vm.Push(tableKinder)
	vm.Push(tableMutter)
	vm.Call(
		2,
		1,
	)

	retval := vm.LState.Get(1)
	vm.Pop(1)

	if retval.Type() != lua.LTNil {
		err = errors.Errorf("lua error: %s", retval)
		return
	}

	if err = sku_fmt.FromLuaTable(kinder, vm.LState, tableKinder); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
