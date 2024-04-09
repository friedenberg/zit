package sku_fmt

import (
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	lua "github.com/yuin/gopher-lua"
)

func Lua(o *sku.Transacted, l *lua.LState, t *lua.LTable) {
	l.SetField(t, "Kennung", lua.LString(o.GetKennung().String()))
	l.SetField(t, "Gattung", lua.LString(o.GetGattung().GetGattungString()))
	l.SetField(t, "Typ", lua.LString(o.GetTyp().String()))

	etiketten := l.NewTable()

	o.Metadatei.GetEtiketten().EachPtr(
		func(e *kennung.Etikett) (err error) {
			l.SetField(etiketten, e.String(), lua.LBool(true))
			return
		},
	)

	l.SetField(t, "Etiketten", etiketten)

	etiketten = l.NewTable()

	o.Metadatei.Verzeichnisse.GetImplicitEtiketten().EachPtr(
		func(e *kennung.Etikett) (err error) {
			l.SetField(etiketten, e.String(), lua.LBool(true))
			return
		},
	)

	l.SetField(t, "EtikettenImplicit", etiketten)
}
