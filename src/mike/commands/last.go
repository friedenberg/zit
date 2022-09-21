package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/delta/transaktion"
	"github.com/friedenberg/zit/src/kilo/umwelt"
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

	var transaktion transaktion.Transaktion

	if transaktion, err = s.ReadLastTransaktion(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, o := range transaktion.Objekten {
		switch o.Gattung {
		case gattung.Zettel:
			errors.PrintOut(o)

		default:
			continue
		}
	}

	return
}
