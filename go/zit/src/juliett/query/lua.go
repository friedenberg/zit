package query

import (
	"io"
	"os"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/src/charlie/files"
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

	ml = &Lua{
		statePool: pool.MakePool(
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
		),
	}

	return
}

type Lua struct {
	statePool schnittstellen.Pool[luaSku, *luaSku]
}

type LuaFlag struct {
	Lua
	value string
}

func (l *LuaFlag) Set(script string) (err error) {
	l.value = script

	var f *os.File

	f, err = files.Open(script)

	if errors.IsNotExist(err) {
		err = nil
	} else if !errors.IsNotExist(err) && err != nil {
		err = errors.Wrap(err)
		return
	} else if err == nil {
		var sb strings.Builder

		if _, err = io.Copy(&sb, f); err != nil {
			err = errors.Wrap(err)
			return
		}

		script = sb.String()

		if err = f.Close(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

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

	l.statePool = pool.MakePool(
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

func (lf *LuaFlag) String() string {
	return lf.value
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

func (*Lua) MatcherLen() int {
	return 0
}

func (*Lua) Each(f schnittstellen.FuncIter[sku.Query]) (err error) {
	return
}
