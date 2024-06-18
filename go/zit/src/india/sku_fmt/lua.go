package sku_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	lua "github.com/yuin/gopher-lua"
)

func ToLuaTable(o *sku.Transacted, l *lua.LState, t *lua.LTable) {
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

func FromLuaTable(o *sku.Transacted, l *lua.LState, t *lua.LTable) (err error) {
	if err = o.Kennung.Set(l.GetField(t, "Kennung").String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	et := l.GetField(t, "Etiketten")
	ets, ok := et.(*lua.LTable)

	if !ok {
		err = errors.Errorf("expected table but got %T", et)
		return
	}

	o.Metadatei.SetEtiketten(nil)

	ets.ForEach(
		func(key, value lua.LValue) {
			var e kennung.Etikett

			if err = e.Set(key.String()); err != nil {
				err = errors.Wrap(err)
				panic(err)
			}

			errors.PanicIfError(o.Metadatei.AddEtikettPtr(&e))
		},
	)

	// TODO Bezeichnung
	// TODO Typ
	// TODO Tai
	// TODO Akte
	// TODO Verzeichnisse

	return
}
