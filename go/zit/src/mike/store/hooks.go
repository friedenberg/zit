package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/hotel/type_blobs"
)

func (s *Store) tryNewHook(
	kinder *sku.Transacted,
	o sku.CommitOptions,
) (err error) {
	if !o.Mode.Contains(object_mode.ModeHooks) {
		return
	}

	var t *sku.Transacted

	if t, err = s.ReadOneObjectId(kinder.GetType()); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	var blob type_blobs.Blob

	if blob, _, err = s.GetBlobStore().GetType().ParseTypedBlob(
		t.GetType(),
		t.GetBlobSha(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer s.GetBlobStore().GetType().PutTypedBlob(t.GetType(), blob)

	script := blob.GetStringLuaHooks()

	if script == "" {
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

	if mutter, err = s.ReadOneObjectId(kinder.GetObjectId()); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	var t *sku.Transacted

	if t, err = s.ReadOneObjectId(kinder.GetType()); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	var blob type_blobs.Blob

	if blob, _, err = s.GetBlobStore().GetType().ParseTypedBlob(
		t.GetType(),
		t.GetBlobSha(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer s.GetBlobStore().GetType().PutTypedBlob(t.GetType(), blob)

	script := blob.GetStringLuaHooks()

	if script == "" {
		return
	}

	if err = s.tryHookWithName(
		kinder,
		mutter,
		sku.CommitOptions{},
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
	o sku.CommitOptions,
) (err error) {
	if !o.Mode.Contains(object_mode.ModeHooks) &&
		!o.Mode.Contains(object_mode.ModeAddToInventoryList) {
		return
	}

	type hook struct {
		script      string
		description string
	}

	hooks := []hook{}

	var t *sku.Transacted

	if t, err = s.ReadOneObjectId(kinder.GetType()); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	var blob type_blobs.Blob

	if blob, _, err = s.GetBlobStore().GetType().ParseTypedBlob(
		t.GetType(),
		t.GetBlobSha(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer s.GetBlobStore().GetType().PutTypedBlob(t.GetType(), blob)

	script := blob.GetStringLuaHooks()

	hooks = append(hooks, hook{script: script, description: "type"})
	hooks = append(hooks, hook{script: s.GetConfig().Hooks, description: "config-mutable"})

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
			err = errors.Wrapf(err, "Type: %q", kinder.GetType())

			if ui.Continue("hook failed", err) {
				err = nil
			} else {
				return
			}
		}
	}

	return
}

func (s *Store) tryPreCommitHook(
	kinder *sku.Transacted,
	mutter *sku.Transacted,
	mode object_mode.Mode,
	selbst *sku.Transacted,
	script string,
) (err error) {
	if mode == object_mode.ModeEmpty {
		return
	}

	var vp sku.LuaVMPoolV1

	if vp, err = s.MakeLuaVMPoolV1(selbst, script); err != nil {
		err = errors.Wrap(err)
		return
	}

	var vm *sku.LuaVMV1

	if vm, err = vp.Get(); err != nil {
		err = errors.Wrap(err)
		return
	}

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

	sku.ToLuaTableV1(
		kinder,
		vm.LState,
		tableKinder,
	)

	var tableMutter *sku.LuaTableV1

	if mutter != nil {
		tableMutter = vm.TablePool.Get()
		defer vm.TablePool.Put(tableMutter)

		sku.ToLuaTableV1(
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

	if err = sku.FromLuaTableV1(kinder, vm.LState, tableKinder); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO add method with hook with reader
func (s *Store) tryHookWithName(
	kinder *sku.Transacted,
	mutter *sku.Transacted,
	o sku.CommitOptions,
	self *sku.Transacted,
	script string,
	name string,
) (err error) {
	var vp sku.LuaVMPoolV1

	if vp, err = s.MakeLuaVMPoolV1(self, script); err != nil {
		err = errors.Wrap(err)
		return
	}

	var vm *sku.LuaVMV1

	if vm, err = vp.Get(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err != nil {
		return
	}

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

	sku.ToLuaTableV1(
		kinder,
		vm.LState,
		tableKinder,
	)

	var tableMutter *sku.LuaTableV1

	if mutter != nil {
		tableMutter = vm.TablePool.Get()
		defer vm.TablePool.Put(tableMutter)

		sku.ToLuaTableV1(
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

	if err = sku.FromLuaTableV1(kinder, vm.LState, tableKinder); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
