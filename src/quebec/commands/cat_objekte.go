package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/delta/sha_collections"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type CatObjekte struct{}

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

func (c CatObjekte) RunWithIds(
	u *umwelt.Umwelt,
	ids kennung.Set,
) (err error) {
	shas := ids.Shas.Copy()
	return c.akten(u, shas)
}

func (c CatObjekte) akten(
	u *umwelt.Umwelt,
	shas sha_collections.Set,
) (err error) {
	// TODO-P3 refactor into reusable
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

	if err = u.Standort().ReadAllShasForGattung(
		gattung.Akte,
		collections.MakeChain(
			shas.WriterContainer(collections.MakeErrStopIteration()),
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
