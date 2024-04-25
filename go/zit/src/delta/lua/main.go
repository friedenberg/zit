package lua

import (
	lua "github.com/yuin/gopher-lua"
)

const (
	LTNil      = lua.LTNil
	LTFunction = lua.LTFunction
	LTTable    = lua.LTTable
	MultRet    = lua.MultRet
)

type (
	LTable        = lua.LTable
	LValue        = lua.LValue
	LState        = lua.LState
	LFunction     = lua.LFunction
	LString       = lua.LString
	FunctionProto = lua.FunctionProto
	LGFunction    = lua.LGFunction
)
