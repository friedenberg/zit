package query

import (
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
	lua "github.com/yuin/gopher-lua"
	lua_ast "github.com/yuin/gopher-lua/ast"
	lua_parse "github.com/yuin/gopher-lua/parse"
)

type luaSku struct {
	*lua.LState
	*lua.LTable
}

func MakeLua(script string) (ml *Lua, err error) {
	ml = &Lua{}

	if err = ml.Set(script); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type Lua struct {
	statePool schnittstellen.Pool[luaSku, *luaSku]
}

func (matcher *Lua) Set(script string) (err error) {
	reader := strings.NewReader(script)

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

	matcher.statePool = pool.MakePool(
		func() (l *luaSku) {
			l = &luaSku{
				LState: lua.NewState(),
			}

			l.LTable = l.NewTable()

			lfunc := l.NewFunctionFromProto(compiled)
			l.Push(lfunc)
			l.PCall(0, lua.MultRet, nil)

			return l
		},
		func(s *luaSku) {
			s.SetTop(0)
		},
	)

	return
}

func (matcher *Lua) ContainsSku(sk *sku.Transacted) bool {
	s := matcher.statePool.Get()
	defer matcher.statePool.Put(s)

	f := s.GetGlobal("contains_matchable")
	s.Push(f)

	sku_fmt.Lua(
		sk,
		s.LState,
		s.LTable,
	)
	s.Push(s.LTable)
	s.Call(
		1,
		1,
	)

	const idx = -1
	return s.CheckBool(idx)
}
