package commands

import (
	"flag"

	"github.com/friedenberg/zit/india/store_with_lock"
)

type RmObjekte struct {
	Type _Type
}

func init() {
	registerCommand(
		"rm-objekte",
		func(f *flag.FlagSet) Command {
			c := &RmObjekte{
				Type: _TypeUnknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithLockedStore{commandWithId{c}}
		},
	)
}

func (c RmObjekte) RunWithId(store store_with_lock.Store, ids ..._Id) (err error) {
	switch c.Type {

	// case _TypeAkte:
	// 	return c.akten(store, ids...)

	case _TypeZettel:
		return c.zettelen(store, ids...)

	default:
		err = _Errorf("unsupported objekte type: %s", c.Type)
		return
	}

	return
}

// func (c RmObjekte) akten(store store_with_lock.Store, ids ..._Id) (err error) {
// 	for _, id := range ids {
// 		var sb _Sha

// 		switch i := id.(type) {
// 		case _Sha:
// 			sb = i

// 		case _Hinweis:
// 			var named _NamedZettel

// 			if named, err = store.Zettels().Read(i); err != nil {
// 				err = _Error(err)
// 				return
// 			}

// 			sb = named.Zettel.Akte

// 		default:
// 			err = _Errorf("unsupported id type: %q", i)
// 			return
// 		}

// 		p := store.DirAkte()

// 		if sb, err = sb.Glob(p); err != nil {
// 			err = _Error(err)
// 			return
// 		}

// 		if err = _ObjekteRead(store.Out, store.Age(), _IdPath(sb, p)); err != nil {
// 			err = _Error(err)
// 			return
// 		}
// 	}

// 	return
// }

func (c RmObjekte) zettelen(store store_with_lock.Store, ids ..._Id) (err error) {
	for _, id := range ids {
		var z _NamedZettel

		if z, err = store.Zettels().Delete(id); err != nil {
			err = _Error(err)
			return
		}

		_Outf("[%s %s] (deleted)\n", z.Hinweis, z.Sha)
	}

	return
}
