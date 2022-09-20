package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/zk_types"
	"github.com/friedenberg/zit/src/charlie/id"
	"github.com/friedenberg/zit/src/delta/id_set"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/kilo/umwelt"
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

			return commandWithIds{
				CommandWithIds: c,
			}
		},
	)
}

func (c CatObjekte) RunWithIds(store *umwelt.Umwelt, ids id_set.Set) (err error) {
	switch c.Type {

	case zk_types.TypeAkte:
		return c.akten(store, ids)

	case zk_types.TypeZettel:
		return c.zettelen(store, ids)

	default:
		err = errors.Errorf("unsupported objekte type: %s", c.Type)
		return
	}
}

func (c CatObjekte) akteShasFromIds(
	u *umwelt.Umwelt,
	ids id_set.Set,
) (shas []sha.Sha, err error) {
	shas = ids.Shas()

	for _, h := range ids.Hinweisen() {
		var zc zettel_checked_out.Zettel

		if zc, err = u.StoreWorkingDirectory().Read(h.String() + ".md"); err != nil {
			err = errors.Wrap(err)
			return
		}

		if zc.State == zettel_checked_out.StateExistsAndDifferent {
			shas = append(shas, zc.External.Named.Stored.Zettel.Akte)
		} else {
			shas = append(shas, zc.Internal.Named.Stored.Zettel.Akte)
		}
	}

	return
}

func (c CatObjekte) akten(store *umwelt.Umwelt, ids id_set.Set) (err error) {
	var shas []sha.Sha

	if shas, err = c.akteShasFromIds(store, ids); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, sb := range shas {
		func(sb sha.Sha) {
			var r io.ReadCloser

			if r, err = store.StoreObjekten().AkteReader(sb); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.PanicIfError(r.Close)

			if io.Copy(store.Out(), r); err != nil {
				err = errors.Wrap(err)
				return
			}
		}(sb)
	}

	return
}

func (c CatObjekte) zettelen(store *umwelt.Umwelt, ids ...id_set.Set) (err error) {
	for _, is := range ids {
		var i id.IdMitKorper
		ok := false

		if i, ok = is.AnyShaOrHinweis(); !ok {
			errors.PrintErrf("unsupported id type: %s", is)
			err = nil
			continue
		}

		var tz zettel_transacted.Zettel

		if tz, err = store.StoreObjekten().ReadOne(i); err != nil {
			err = errors.Wrap(err)
			return
		}

		f := zettel.Objekte{}

		errors.PrintDebug(tz)

		if _, err = f.WriteTo(tz.Named.Stored.Zettel, store.Out()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
