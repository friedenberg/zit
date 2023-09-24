package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/sha_collections"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type CatObjekte struct{}

func init() {
	registerCommand(
		"cat-objekte",
		func(f *flag.FlagSet) Command {
			c := &CatObjekte{}

			return c
		},
	)
}

func (c CatObjekte) Run(
	u *umwelt.Umwelt,
	args ...string,
) (err error) {
	shas := collections_value.MakeMutableValueSet[sha.Sha](nil)

	for _, v := range args {
		if err = iter.AddString[sha.Sha](
			shas,
			v,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return c.akten(u, shas)
}

func (c CatObjekte) akten(
	u *umwelt.Umwelt,
	shas sha_collections.Set,
) (err error) {
	// TODO-P3 refactor into reusable
	akteWriter := iter.MakeSyncSerializer(
		func(rc io.ReadCloser) (err error) {
			defer errors.DeferredCloser(&err, rc)

			if _, err = io.Copy(u.Out(), rc); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)

	if err = u.Standort().ReadAllShasForGattung(
		u.Konfig().GetStoreVersion(),
		gattung.Akte,
		iter.MakeChain(
			collections.WriterContainer(shas, collections.MakeErrStopIteration()),
			func(sb sha.Sha) (err error) {
				var r io.ReadCloser

				if r, err = u.Standort().AkteReader(sb); err != nil {
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
