package sku_fmt

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	lua "github.com/yuin/gopher-lua"
)

func Lua(o *sku.Transacted, ki kennung.Index, l *lua.LState, t *lua.LTable) {
	l.SetField(t, "Kennung", lua.LString(o.GetKennungLike().String()))
	l.SetField(t, "Gattung", lua.LString(o.GetGattung().GetGattungString()))
	l.SetField(t, "Typ", lua.LString(o.GetTyp().String()))

	etiketten := l.NewTable()

	addOne := func(e *kennung.Etikett) (err error) {
		l.SetField(etiketten, e.String(), lua.LBool(true))
		return
	}

	o.Metadatei.GetEtiketten().EachPtr(
		func(e *kennung.Etikett) (err error) {
			// indexed, err := ki.Etiketten(e)

			// if err == nil {
			// 	indexed.GetExpandedRight().EachPtr(addOne)
			// }

			return addOne(e)
		},
	)

	l.SetField(t, "Etiketten", etiketten)

	return
}
