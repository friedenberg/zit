package commands

import (
	"flag"

	"github.com/friedenberg/zit/india/store_with_lock"
)

type CatObjekte struct {
	Type _Type
}

func init() {
	registerCommand(
		"cat-objekte",
		func(f *flag.FlagSet) Command {
			c := &CatObjekte{
				Type: _TypeUnknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithLockedStore{commandWithId{c}}
		},
	)
}

func (c CatObjekte) RunWithId(store store_with_lock.Store, ids ..._Id) (err error) {
	switch c.Type {

	case _TypeAkte:
		return c.akten(store, ids...)

	case _TypeZettel:
		return c.zettelen(store, ids...)

	default:
		err = _Errorf("unsupported objekte type: %s", c.Type)
		return
	}

	return
}

func (c CatObjekte) akten(store store_with_lock.Store, ids ..._Id) (err error) {
	for _, id := range ids {
		var sb _Sha

		switch i := id.(type) {
		case _Sha:
			sb = i

		case _Hinweis:
			var named _NamedZettel

			if named, err = store.Zettels().Read(i); err != nil {
				err = _Error(err)
				return
			}

			sb = named.Zettel.Akte

		default:
			err = _Errorf("unsupported id type: %q", i)
			return
		}

		p := store.DirAkte()

		if sb, err = sb.Glob(p); err != nil {
			err = _Error(err)
			return
		}

		if err = _ObjekteRead(store.Out, store.Age(), _IdPath(sb, p)); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}

func (c CatObjekte) zettelen(store store_with_lock.Store, ids ..._Id) (err error) {
	for _, id := range ids {
		var z _NamedZettel

		if z, err = store.Zettels().Read(id); err != nil {
			err = _Error(err)
			return
		}

		f := _StoredZettelFormatsObjekte{}

		if _, err = f.WriteTo(z.Stored, store.Out); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}
