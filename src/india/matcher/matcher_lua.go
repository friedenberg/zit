package matcher

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/kennung"
	lua "github.com/yuin/gopher-lua"
)

func MakeMatcherWithLua(base, conditional kennung.Etikett, str string) Matcher {
	l := lua.NewState()

	l.DoString(
		`function f (base, conditional, sku) return base .. conditional end`,
	)

	return &matcherWithLua{
		base:        base,
		conditional: conditional,
		state:       l,
	}
}

type matcherWithLua struct {
	base        kennung.Etikett
	conditional kennung.Etikett
	state       *lua.LState
}

func (m matcherWithLua) String() string {
	return "lua"
	// sb := &strings.Builder{}

	// if m.Matcher != nil {
	// 	sb.WriteString(m.Matcher.String())
	// }

	// sb.WriteString(m.Sigil.String())

	// return sb.String()
}

func (matcher matcherWithLua) ContainsMatchable(matchable Matchable) bool {
	if matcher.state == nil {
		return true
	}

	f := matcher.state.GetGlobal("f")
	matcher.state.Push(f)
	matcher.state.Push(lua.LString("wow"))
	matcher.state.Push(lua.LString("ok"))
	// matcher.state.Push(lua.LString(matcher.base.String()))
	// matcher.state.Push(lua.LString(matcher.conditional.String()))
	matcher.state.Push(lua.LString("womp"))
	matcher.state.Call(
		3,
		1,
	) // Call f with 0 arguments and 1 result.

	const idx = -1
	return true
	// return matcher.state.CheckBool(idx)
}

func (_ matcherWithLua) MatcherLen() int {
	return 0
}

func (_ matcherWithLua) Each(f schnittstellen.FuncIter[Matcher]) (err error) {
	return
}
