package query

import (
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/src/echo/kennung"
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

func MakeLua(ki kennung.Index, script string) (m Matcher, err error) {
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

	m = &matcherLua{
		kennungIndex: ki,
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

type matcherLua struct {
	kennungIndex kennung.Index
	statePool    schnittstellen.Pool[luaSku, *luaSku]
}

func (m *matcherLua) String() string {
	return "lua"
	// sb := &strings.Builder{}

	// if m.Matcher != nil {
	// 	sb.WriteString(m.Matcher.String())
	// }

	// sb.WriteString(m.Sigil.String())

	// return sb.String()
}

func (matcher *matcherLua) ContainsMatchable(matchable *sku.Transacted) bool {
	s := matcher.statePool.Get()
	defer matcher.statePool.Put(s)

	f := s.GetGlobal("contains_matchable")
	s.Push(f)

	sku_fmt.Lua(
		matchable,
		matcher.kennungIndex,
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

func (*matcherLua) MatcherLen() int {
	return 0
}

func (_ *matcherLua) Each(f schnittstellen.FuncIter[Matcher]) (err error) {
	return
}
