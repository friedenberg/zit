package sku_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO setup versions of LuaTable
type LuaTable struct {
	Transacted   *lua.LTable
	Tags         *lua.LTable
	TagsImplicit *lua.LTable
}

func ToLuaTable(tg sku.TransactedGetter, l *lua.LState, t *LuaTable) {
	o := tg.GetSku()

	l.SetField(t.Transacted, "Gattung", lua.LString(o.GetGenre().String()))
	l.SetField(t.Transacted, "Kennung", lua.LString(o.GetObjectId().String()))
	l.SetField(t.Transacted, "Typ", lua.LString(o.GetType().String()))

	tags := t.Tags

	o.Metadata.GetTags().EachPtr(
		func(e *ids.Tag) (err error) {
			l.SetField(tags, e.String(), lua.LBool(true))
			return
		},
	)

	tags = t.TagsImplicit

	o.Metadata.Cache.GetImplicitTags().EachPtr(
		func(e *ids.Tag) (err error) {
			l.SetField(tags, e.String(), lua.LBool(true))
			return
		},
	)
}

func FromLuaTable(o *sku.Transacted, l *lua.LState, lt *LuaTable) (err error) {
	t := lt.Transacted

	g := genres.MakeOrUnknown(l.GetField(t, "Gattung").String())

	o.ObjectId.SetGenre(g)
	k := l.GetField(t, "Kennung").String()

	if k != "" {
		if err = o.ObjectId.Set(k); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	et := l.GetField(t, "Etiketten")
	ets, ok := et.(*lua.LTable)

	if !ok {
		err = errors.Errorf("expected table but got %T", et)
		return
	}

	o.Metadata.SetTags(nil)

	ets.ForEach(
		func(key, value lua.LValue) {
			var e ids.Tag

			if err = e.Set(key.String()); err != nil {
				err = errors.Wrap(err)
				panic(err)
			}

			errors.PanicIfError(o.Metadata.AddTagPtr(&e))
		},
	)

	// TODO Bezeichnung
	// TODO Typ
	// TODO Tai
	// TODO Blob
	// TODO Verzeichnisse

	return
}
