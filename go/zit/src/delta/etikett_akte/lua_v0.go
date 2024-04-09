package etikett_akte

var LuaV0Typ = "etikett-lua"

type LuaV0 string

func (l *LuaV0) Reset() {
	*l = ""
}

func (a *LuaV0) ResetWith(b *LuaV0) {
	*a = *b
}
