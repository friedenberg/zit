package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/charlie/zk_types"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/golf/zettel_formats"
	"github.com/friedenberg/zit/src/india/store_with_lock"
)

type CatObjekte struct {
	Type zk_types.Type
}

func init() {
	registerCommand(
		"cat-objekte",
		func(f *flag.FlagSet) Command {
			c := &CatObjekte{
				Type: zk_types.TypeUnknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithLockedStore{commandWithId{c}}
		},
	)
}

func (c CatObjekte) RunWithId(store store_with_lock.Store, ids ...id.Id) (err error) {
	switch c.Type {

	case zk_types.TypeAkte:
		return c.akten(store, ids...)

	case zk_types.TypeZettel:
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
			var tz stored_zettel.Transacted

			if tz, err = store.StoreObjekten().Read(i); err != nil {
				err = errors.Error(err)
				return
			}

			sb = tz.Named.Stored.Zettel.Akte

		default:
			err = errors.Errorf("unsupported id type: %q", i)
			return
		}

		func(sb sha.Sha) {
			var r io.ReadCloser

			if r, err = store.StoreObjekten().AkteReader(sb); err != nil {
				err = errors.Error(err)
				return
			}

			defer stdprinter.PanicIfError(r.Close)

			if io.Copy(store.Out, r); err != nil {
				err = errors.Error(err)
				return
			}
		}(sb)
	}

	return
}

func (c CatObjekte) zettelen(store store_with_lock.Store, ids ...id.Id) (err error) {
	for _, id := range ids {
		var tz stored_zettel.Transacted

		if tz, err = store.StoreObjekten().Read(id); err != nil {
			err = errors.Error(err)
			return
		}

		f := zettel_formats.Objekte{}

		logz.PrintDebug(tz)

		if _, err = f.WriteTo(tz.Named.Stored.Zettel, store.Out); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}
