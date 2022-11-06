package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/id_set"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
	"github.com/friedenberg/zit/src/mike/umwelt"
)

type CatObjekte struct {
	Type gattung.Gattung
}

func init() {
	registerCommand(
		"cat-objekte",
		func(f *flag.FlagSet) Command {
			c := &CatObjekte{
				Type: gattung.Unknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithIds{
				CommandWithIds: c,
			}
		},
	)
}

func (c CatObjekte) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	is = id_set.MakeProtoIdSet(
		id_set.ProtoId{
			MutableId: &sha.Sha{},
		},
	)

	return
}

func (c CatObjekte) RunWithIds(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	shas := ids.Shas()

	switch c.Type {
	case gattung.Akte:
		return c.akten(u, shas)

	case gattung.Zettel:
		return c.zettelen(u, shas)

	default:
		err = errors.Errorf("unsupported objekte type: %s", c.Type)
		return
	}
}

func (c CatObjekte) akten(u *umwelt.Umwelt, shas sha.Set) (err error) {
	akteWriter := collections.MakeSyncSerializer(
		func(rc io.ReadCloser) (err error) {
			defer errors.Deferred(&err, rc.Close)

			if _, err = io.Copy(u.Out(), rc); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)

	if err = u.StoreObjekten().ReadAllAktenShas(
		collections.MakeChain(
			shas.WriterContainer(),
			func(sb sha.Sha) (err error) {
				var r io.ReadCloser

				if r, err = u.StoreObjekten().AkteReader(sb); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = akteWriter(r); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c CatObjekte) zettelen(u *umwelt.Umwelt, shas sha.Set) (err error) {
	w := collections.MakeChain(
		zettel_transacted.MakeWriterZettelNamed(
			func(z *zettel_named.Zettel) (err error) {
				if !shas.Contains(z.Stored.Sha) {
					err = io.EOF
				}

				return
			},
		),
		zettel_transacted.MakeWriterZettel(
			zettel.MakeSerializedFormatWriter(
				&zettel.Objekte{},
				u.Out(),
				u.StoreObjekten(),
				u.Konfig(),
			),
		),
	)

	if err = u.StoreObjekten().ReadAllSchwanzenTransacted(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
