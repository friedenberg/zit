package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/delta/typ_akte"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

func (s *Store) tryNewHook(
	kinder *sku.Transacted,
	o ObjekteOptions,
) (err error) {
	if !o.Mode.Contains(objekte_mode.ModeHooks) {
		return
	}

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

	if err = s.tryHookWithName(
		kinder,
		nil,
		o,
		t,
		script,
		"on_new",
	); err != nil {
		err = errors.Wrapf(err, "Hook: %#v", script)
		return
	}

	return
}

func (s *Store) TryFormatHook(
	kinder *sku.Transacted,
) (err error) {
	var mutter *sku.Transacted

	if mutter, err = s.ReadOneKennung(kinder.GetKennung()); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

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

	if err = s.tryHookWithName(
		kinder,
		mutter,
		ObjekteOptions{},
		t,
		script,
		"on_format",
	); err != nil {
		err = errors.Wrapf(err, "Hook: %#v", script)
		return
	}

	return
}

func (s *Store) tryPreCommitHooks(
	kinder *sku.Transacted,
	mutter *sku.Transacted,
	o ObjekteOptions,
) (err error) {
	if !o.Mode.Contains(objekte_mode.ModeHooks) &&
		!o.Mode.Contains(objekte_mode.ModeAddToBestandsaufnahme) {
		return
	}

	type hook struct {
		script      string
		description string
	}

	hooks := []hook{}

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

	script, _ := akte.Hooks.(string)

	hooks = append(hooks, hook{script: script, description: "typ"})
	hooks = append(hooks, hook{script: s.GetKonfig().Hooks, description: "erworben"})

	for _, h := range hooks {
		if h.script == "" {
			continue
		}

		if err = s.tryHookWithName(
			kinder,
			mutter,
			o,
			t,
			h.script,
			"on_pre_commit",
		); err != nil {
			err = errors.Wrapf(err, "Hook: %#v", h)
			return
		}
	}

	return
}

func (s *Store) tryPreCommitHook(
	kinder *sku.Transacted,
	mutter *sku.Transacted,
	mode objekte_mode.Mode,
	selbst *sku.Transacted,
	script string,
) (err error) {
	if mode == objekte_mode.ModeEmpty {
		return
	}

	var vp sku_fmt.LuaVMPool

	if vp, err = s.MakeLuaVMPool(selbst, script); err != nil {
		err = errors.Wrap(err)
		return
	}

	vm := vp.Get()
	defer vp.Put(vm)

	var tt *lua.LTable

	if tt, err = vm.GetTopTableOrError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	f := vm.GetField(tt, "on_pre_commit")

	if f.Type() != lua.LTFunction {
		return
	}

	tableKinder := vm.TablePool.Get()
	defer vm.TablePool.Put(tableKinder)

	sku_fmt.ToLuaTable(
		kinder,
		vm.LState,
		tableKinder,
	)

	var tableMutter *sku_fmt.LuaTable

	if mutter != nil {
		tableMutter = vm.TablePool.Get()
		defer vm.TablePool.Put(tableMutter)

		sku_fmt.ToLuaTable(
			mutter,
			vm.LState,
			tableMutter,
		)
	}

	vm.Push(f)
	vm.Push(tableKinder.Transacted)
	vm.Push(tableMutter.Transacted)

	if err = vm.PCall(2, 1, nil); err != nil {
		err = errors.Wrap(err)
		return
	}

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

func (s *Store) tryHookWithName(
	kinder *sku.Transacted,
	mutter *sku.Transacted,
	o ObjekteOptions,
	selbst *sku.Transacted,
	script string,
	name string,
) (err error) {
	var vp sku_fmt.LuaVMPool

	if vp, err = s.MakeLuaVMPool(selbst, script); err != nil {
		err = errors.Wrap(err)
		return
	}

	vm := vp.Get()
	defer vp.Put(vm)

	var tt *lua.LTable

	if tt, err = vm.GetTopTableOrError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	f := vm.GetField(tt, name)

	if f.Type() != lua.LTFunction {
		return
	}

	tableKinder := vm.TablePool.Get()
	defer vm.TablePool.Put(tableKinder)

	sku_fmt.ToLuaTable(
		kinder,
		vm.LState,
		tableKinder,
	)

	var tableMutter *sku_fmt.LuaTable

	if mutter != nil {
		tableMutter = vm.TablePool.Get()
		defer vm.TablePool.Put(tableMutter)

		sku_fmt.ToLuaTable(
			mutter,
			vm.LState,
			tableMutter,
		)
	}

	vm.Push(f)
	vm.Push(tableKinder.Transacted)

	if tableMutter != nil {
		vm.Push(tableMutter.Transacted)
	} else {
		vm.Push(nil)
	}

	if err = vm.PCall(2, 1, nil); err != nil {
		err = errors.Wrap(err)
		return
	}

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
