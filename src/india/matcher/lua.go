package matcher

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/pool"
	"github.com/friedenberg/zit/src/india/sku_formats"
	lua "github.com/yuin/gopher-lua"
)

func MakeMatcherLua(script string) Matcher {
	return &matcherLua{
		statePool: pool.MakePool(
			func() *lua.LState {
				l := lua.NewState()

				l.DoString(
					script,
				)

				return l
			},
			func(l *lua.LState) {
				l.SetTop(0)
			},
		),
	}
}

type matcherLua struct {
	lock      sync.Mutex
	statePool schnittstellen.Pool[lua.LState, *lua.LState]
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

func (matcher *matcherLua) ContainsMatchable(matchable Matchable) bool {
	s := matcher.statePool.Get()
	defer matcher.statePool.Put(s)

	f := s.GetGlobal("contains_matchable")
	s.Push(f)
	s.Push(lua.LString(sku_formats.String(matchable.GetSkuLike())))
	s.Call(
		1,
		1,
	)

	const idx = -1
	return s.CheckBool(idx)
}

func (_ *matcherLua) MatcherLen() int {
	return 0
}

func (_ *matcherLua) Each(f schnittstellen.FuncIter[Matcher]) (err error) {
	return
}
