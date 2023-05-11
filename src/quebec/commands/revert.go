package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/delta/kennung"
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

func (c Revert) RunWithIds(u *umwelt.Umwelt, ids kennung.Set) (err error) {
	return todo.Implement()
	// timestamps := ids.Timestamps.ImmutableClone()
	// var transaktion *transaktion.Transaktion

	// s := u.StoreObjekten()

	// if timestamps.Len() == 1 {
	// 	if transaktion, err = s.GetTransaktionStore().ReadTransaktion(
	// 		timestamps.Any(),
	// 	); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}
	// } else {
	// 	if transaktion, err = s.GetTransaktionStore().ReadLastTransaktion(); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// 	errors.Out().Printf(
	// 		"ignoring arguments and using last transkation: %s",
	// 		transaktion.Time,
	// 	)
	// }

	// if err = u.Lock(); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	// defer errors.Deferred(&err, u.Unlock)

	// var zts zettel.MutableSet

	// if zts, err = s.RevertTransaktion(transaktion); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	// zts.Each(u.PrinterZettelTransactedDelta())

	// return
}
