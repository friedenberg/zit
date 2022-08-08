package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/node_type"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/objekte"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/golf/stored_zettel_formats"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type CatObjekte struct {
	Type node_type.Type
}

func init() {
	registerCommand(
		"cat-objekte",
		func(f *flag.FlagSet) Command {
			c := &CatObjekte{
				Type: node_type.TypeUnknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithLockedStore{commandWithId{c}}
		},
	)
}

func (c CatObjekte) RunWithId(store store_with_lock.Store, ids ...id.Id) (err error) {
	switch c.Type {

	case node_type.TypeAkte:
		return c.akten(store, ids...)

	case node_type.TypeZettel:
		return c.zettelen(store, ids...)

	default:
		err = errors.Errorf("unsupported objekte type: %s", c.Type)
		return
	}
}

func (c CatObjekte) akten(store store_with_lock.Store, ids ...id.Id) (err error) {
	for _, idt := range ids {
		var sb sha.Sha

		switch i := idt.(type) {
		case sha.Sha:
			sb = i

		case hinweis.Hinweis:
			var named stored_zettel.Named

			if named, err = store.Zettels().Read(i); err != nil {
				err = errors.Error(err)
				return
			}

			sb = named.Zettel.Akte

		default:
			err = errors.Errorf("unsupported id type: %q", i)
			return
		}

		p := store.DirAkte()

		if sb, err = sb.Glob(p); err != nil {
			err = errors.Error(err)
			return
		}

		if err = objekte.Read(store.Out, store.Age(), id.Path(sb, p)); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}

func (c CatObjekte) zettelen(store store_with_lock.Store, ids ...id.Id) (err error) {
	for _, id := range ids {
		var z stored_zettel.Named

		if z, err = store.Zettels().Read(id); err != nil {
			err = errors.Error(err)
			return
		}

		f := stored_zettel_formats.Objekte{}

		if _, err = f.WriteTo(z.Stored, store.Out); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}
