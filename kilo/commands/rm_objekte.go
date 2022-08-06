package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/node_type"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type RmObjekte struct {
	Type node_type.Type
}

func init() {
	registerCommand(
		"rm-objekte",
		func(f *flag.FlagSet) Command {
			c := &RmObjekte{
				Type: node_type.TypeUnknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithLockedStore{commandWithId{c}}
		},
	)
}

func (c RmObjekte) RunWithId(store store_with_lock.Store, ids ...id.Id) (err error) {
	switch c.Type {

	// case node_type.TypeAkte:
	// 	return c.akten(store, ids...)

	case node_type.TypeZettel:
		return c.zettelen(store, ids...)

	default:
		err = errors.Errorf("unsupported objekte type: %s", c.Type)
		return
	}
}

// func (c RmObjekte) akten(store store_with_lock.Store, ids ...id.Id) (err error) {
// 	for _, id := range ids {
// 		var sb sha.Sha

// 		switch i := id.(type) {
// 		case sha.Sha:
// 			sb = i

// 		case hinweis.Hinweis:
// 			var named stored_zettel.Named

// 			if named, err = store.Zettels().Read(i); err != nil {
// 				err = errors.Error(err)
// 				return
// 			}

// 			sb = named.Zettel.Akte

// 		default:
// 			err = errors.Errorf("unsupported id type: %q", i)
// 			return
// 		}

// 		p := store.DirAkte()

// 		if sb, err = sb.Glob(p); err != nil {
// 			err = errors.Error(err)
// 			return
// 		}

// 		if err = objekte.Read(store.Out, store.Age(), id.Path(sb, p)); err != nil {
// 			err = errors.Error(err)
// 			return
// 		}
// 	}

// 	return
// }

func (c RmObjekte) zettelen(store store_with_lock.Store, ids ...id.Id) (err error) {
	for _, id := range ids {
		var z stored_zettel.Named

		if z, err = store.Zettels().Delete(id); err != nil {
			err = errors.Error(err)
			return
		}

		stdprinter.Outf("[%s %s] (deleted)\n", z.Hinweis, z.Sha)
	}

	return
}
