package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type CatObjekte struct {
}

func init() {
	registerCommand(
		"cat-objekte",
		func(f *flag.FlagSet) Command {
			c := &CatObjekte{}

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
			Expand: func(v string) (out string, err error) {
				var s sha.Sha
				s, err = u.StoreObjekten().Abbr().ExpandShaString(v)
				out = s.String()
				return
			},
		},
	)

	return
}

func (c CatObjekte) RunWithIds(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	shas := ids.Shas.Copy()
	return c.akten(u, shas)
}

func (c CatObjekte) akten(u *umwelt.Umwelt, shas sha.Set) (err error) {
	//TODO-P3 refactor into reusable
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
			shas.WriterContainer(io.EOF),
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
