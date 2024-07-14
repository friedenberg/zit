package store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

func (s *Store) MakeLuaVMPoolWithSku(
	sk *sku.Transacted,
) (lvp sku_fmt.LuaVMPool, err error) {
	if sk.GetType().String() != "lua" {
		err = errors.Errorf("unsupported typ: %s, Sku: %s", sk.GetType(), sk)
		return
	}

	var ar sha.ReadCloser

	if ar, err = s.GetStandort().BlobReader(sk.GetBlobSha()); err != nil {
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
) (vp sku_fmt.LuaVMPool, err error) {
	b := u.luaVMPoolBuilder.Clone().
		WithScript(script).
		WithApply(query.MakeSelfApply(selbst))

	var lvmp *lua.VMPool

	if lvmp, err = b.Build(); err != nil {
		err = errors.Wrap(err)
		return
	}

	vp = sku_fmt.MakeLuaVMPool(lvmp, selbst)

	return
}

func (u *Store) MakeLuaVMPoolWithReader(
	selbst *sku.Transacted,
	r io.Reader,
) (vp sku_fmt.LuaVMPool, err error) {
	b := u.luaVMPoolBuilder.Clone().
		WithReader(r).
		WithApply(query.MakeSelfApply(selbst))

	var lvmp *lua.VMPool

	if lvmp, err = b.Build(); err != nil {
		err = errors.Wrap(err)
		return
	}

	vp = sku_fmt.MakeLuaVMPool(lvmp, selbst)

	return
}

// type LuaVMPool struct {
// 	*sku.Transacted
// 	*lua.VMPool
// }

// func (lvm LuaVMPool) Put(vm LuaVM) {
// 	// vm.Put(vm.LTable)
// 	lvm.VMPool.Put(vm.VM)
// }

// func (lvm LuaVMPool) Get() (vm LuaVM, err error) {
// 	vm.VM = lvm.VMPool.Get()

// 	if vm.LValue, err = vm.GetTopTableOrError(); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	selbstTable := vm.Pool.Get()
// 	sku_fmt.ToLuaTable(lvm.Transacted, vm.LState, selbstTable)
// 	vm.SetField(vm.LValue, "Selbst", selbstTable)

// 	return
// }

// type LuaVM struct {
// 	*lua.VM
// 	lua.LValue
// }
