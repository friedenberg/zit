package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Last struct {
	Type gattung.Gattung
}

func init() {
	registerCommand(
		"last",
		func(f *flag.FlagSet) Command {
			c := &Last{}

			return c
		},
	)
}

func (c Last) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) != 0 {
		errors.PrintOut("ignoring arguments")
	}

	s := u.StoreObjekten()

	var transaktion *transaktion.Transaktion

	if transaktion, err = s.ReadLastTransaktion(); err != nil {
		err = errors.Wrap(err)
		return
	}

	transaktion.Each(
		func(o *sku.Sku) (err error) {
			switch o.Gattung {
			case gattung.Zettel:
				errors.PrintOut(o)
			}

			return
		},
	)

	return
}
