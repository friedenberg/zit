package sku_fmt

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	lua "github.com/yuin/gopher-lua"
)

func Lua(o *sku.Transacted, ki kennung.Index, l *lua.LState, t *lua.LTable) {
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
