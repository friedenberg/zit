package sku_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type LuaTable struct {
	Transacted        *lua.LTable
	Etiketten         *lua.LTable
	EtikettenImplicit *lua.LTable
}

func ToLuaTable(o *sku.Transacted, l *lua.LState, t *LuaTable) {
	l.SetField(t.Transacted, "Gattung", lua.LString(o.GetGenre().String()))
	l.SetField(t.Transacted, "Kennung", lua.LString(o.GetKennung().String()))
	l.SetField(t.Transacted, "Gattung", lua.LString(o.GetGenre().GetGenreString()))
	l.SetField(t.Transacted, "Typ", lua.LString(o.GetTyp().String()))

	etiketten := t.Etiketten

	o.Metadatei.GetEtiketten().EachPtr(
		func(e *ids.Tag) (err error) {
			l.SetField(etiketten, e.String(), lua.LBool(true))
			return
		},
	)

	etiketten = t.EtikettenImplicit

	o.Metadatei.Verzeichnisse.GetImplicitEtiketten().EachPtr(
		func(e *ids.Tag) (err error) {
			l.SetField(etiketten, e.String(), lua.LBool(true))
			return
		},
	)
}

func FromLuaTable(o *sku.Transacted, l *lua.LState, lt *LuaTable) (err error) {
	t := lt.Transacted

	var g genres.Genre
	if err = g.Set(l.GetField(t, "Gattung").String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	o.Kennung.SetGenre(g)
	k := l.GetField(t, "Kennung").String()

	if err = o.Kennung.Set(k); err != nil {
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
			var e ids.Tag

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
