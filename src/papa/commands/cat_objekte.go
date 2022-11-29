package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/id_set"
	"github.com/friedenberg/zit/src/india/zettel"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type CatObjekte struct {
	Gattung gattung.Gattung
}

func init() {
	registerCommand(
		"cat-objekte",
		func(f *flag.FlagSet) Command {
			c := &CatObjekte{
				Gattung: gattung.Unknown,
			}

			f.Var(&c.Gattung, "gattung", "ObjekteType")

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

	switch c.Gattung {
	case gattung.Akte:
		return c.akten(u, shas)

	case gattung.Zettel:
		return c.zettelen(u, shas)

		//TODO add typ

	default:
		err = errors.Errorf("unsupported objekte type: %s", c.Gattung)
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
		zettel.MakeWriterZettelNamed(
			func(z *zettel.Named) (err error) {
				if !shas.Contains(z.Stored.Sha) {
					err = io.EOF
				}

				return
			},
		),
		zettel.MakeWriterZettel(
			zettel.MakeSerializedFormatWriter(
				&zettel.Objekte{},
				u.Out(),
				u.StoreObjekten(),
				u.Konfig(),
			),
		),
	)

	if err = u.StoreObjekten().Zettel().ReadAllSchwanzenTransacted(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
