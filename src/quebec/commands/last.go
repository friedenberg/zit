package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/india/typ"
	"github.com/friedenberg/zit/src/oscar/umwelt"
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
		errors.Err().Print("ignoring arguments")
	}

	s := u.StoreObjekten()

	var transaktion *transaktion.Transaktion

	if transaktion, err = s.ReadLastTransaktion(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = transaktion.Skus.Each(
		func(o sku.SkuLike) (err error) {
			switch o.GetGattung() {
			case gattung.Typ:
				var te *typ.Transacted

				if te, err = u.StoreObjekten().Typ().Inflate(
					transaktion.Time,
					o,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				u.PrinterTypTransacted("test")(te)

			default:
				errors.Out().Print(o)
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
