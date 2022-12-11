package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Revert struct {
	Type gattung.Gattung
}

func init() {
	registerCommand(
		"revert",
		func(f *flag.FlagSet) Command {
			c := &Revert{}

			return commandWithIds{
				CommandWithIds: c,
			}
		},
	)
}

func (c Revert) RunWithIds(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	timestamps := ids.Timestamps()
	var transaktion *transaktion.Transaktion

	s := u.StoreObjekten()

	if timestamps.Len() == 1 {
		if transaktion, err = s.ReadTransaktion(timestamps.Any()); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if transaktion, err = s.ReadLastTransaktion(); err != nil {
			err = errors.Wrap(err)
			return
		}

		errors.Out().Printf(
			"ignoring arguments and using last transkation: %s",
			transaktion,
		)
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	var zts zettel.MutableSet

	if zts, err = s.RevertTransaktion(transaktion); err != nil {
		err = errors.Wrap(err)
		return
	}

	zts.Each(u.PrinterZettelTransacted(format.StringUpdated))

	return
}
