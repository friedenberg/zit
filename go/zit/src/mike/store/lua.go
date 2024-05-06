package store

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/lua"
)

func (u *Store) MakeLuaVMPool(script string) (vp *lua.VMPool, err error) {
	vp = u.luaVMPoolBuilder.Build()

	if err = vp.Set(script); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
