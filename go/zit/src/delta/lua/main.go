package lua

import (
	lua "github.com/yuin/gopher-lua"
)

const (
	LTNil      = lua.LTNil
	LTFunction = lua.LTFunction
	LTTable    = lua.LTTable
)

type (
	LTable    = lua.LTable
	LValue    = lua.LValue
	LFunction = lua.LFunction
)
