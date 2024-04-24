package query

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/lua"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
)

func MakeLua(script string) (ml *Lua, err error) {
	ml = &Lua{}

	if err = ml.Set(script); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type Lua struct {
	lua.VMPool
}

func (matcher *Lua) ContainsSku(sk *sku.Transacted) bool {
	vm := matcher.Get()
	defer matcher.Put(vm)

	t := vm.Pool.Get()
	defer vm.Put(t)

	tt, err := vm.GetTopTableOrError()
	errors.PanicIfError(err)

	f := vm.GetField(tt, "contains_sku")

	if f == nil {
		return false
	}

	vm.Push(f)

	sku_fmt.ToLuaTable(
		sk,
		vm.LState,
		t,
	)
	vm.Push(t)
	vm.Call(
		1,
		1,
	)

	// TODO make safer than checking bool
	const idx = -1
	return vm.CheckBool(idx)
}
