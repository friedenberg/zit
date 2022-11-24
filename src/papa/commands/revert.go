package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/echo/id_set"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/november/umwelt"
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

	if len(timestamps) == 1 {
		if transaktion, err = s.ReadTransaktion(timestamps[0]); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if transaktion, err = s.ReadLastTransaktion(); err != nil {
			err = errors.Wrap(err)
			return
		}

		errors.PrintOutf(
			"ignoring arguments and using last transkation: %s",
			transaktion,
		)
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer u.Unlock()

	var zts zettel_transacted.MutableSet

	if zts, err = s.RevertTransaktion(transaktion); err != nil {
		err = errors.Wrap(err)
		return
	}

	zts.Each(u.PrinterZettelTransacted(format.StringUpdated))

	return
}
