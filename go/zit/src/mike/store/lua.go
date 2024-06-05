package store

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/lua"
	"code.linenisgreat.com/zit/src/delta/sha"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
)

func (s *Store) MakeLuaVMPoolWithSku(
	sk *sku.Transacted,
) (lvp LuaVMPool, err error) {
	if sk.GetTyp().String() != "lua" {
		err = errors.Errorf("unsupported etikett: %s", sk)
		return
	}

	var ar sha.ReadCloser

	if ar, err = s.GetStandort().AkteReader(sk.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	if lvp, err = s.MakeLuaVMPoolWithReader(sk, ar); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (u *Store) MakeLuaVMPool(
	selbst *sku.Transacted,
	script string,
) (vp LuaVMPool, err error) {
	vp.Transacted = selbst
	vp.VMPool = u.luaVMPoolBuilder.Build()

	if err = vp.Set(script); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (u *Store) MakeLuaVMPoolWithReader(
	selbst *sku.Transacted,
	r io.Reader,
) (vp LuaVMPool, err error) {
	vp.Transacted = selbst
	vp.VMPool = u.luaVMPoolBuilder.Build()

	if err = vp.SetReader(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type LuaVMPool struct {
	*sku.Transacted
	*lua.VMPool
}

func (lvm LuaVMPool) Put(vm LuaVM) {
	vm.Put(vm.LTable)
	lvm.VMPool.Put(vm.VM)
}

func (lvm LuaVMPool) Get() (vm LuaVM, err error) {
	vm.VM = lvm.VMPool.Get()

	if vm.LTable, err = vm.GetTopTableOrError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	selbstTable := vm.Pool.Get()
	sku_fmt.ToLuaTable(lvm.Transacted, vm.LState, selbstTable)
	vm.SetField(vm.LTable, "Selbst", selbstTable)

	return
}

type LuaVM struct {
	*lua.VM
	*lua.LTable
}
