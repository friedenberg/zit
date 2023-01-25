package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/india/bestandsaufnahme"
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
		errors.Err().Print("ignoring arguments")
	}

	method := c.runWithTransaktion

	if u.Konfig().UseBestandsaufnahme {
		method = c.runWithBestandsaufnahm
	}

	if err = method(u); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Last) runWithBestandsaufnahm(u *umwelt.Umwelt) (err error) {
	s := u.StoreObjekten()

	var b *bestandsaufnahme.Objekte

	if b, err = s.GetBestandsaufnahmeStore().ReadLast(); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.TodoP3("support log line format for skus")
	if err = b.Akte.Skus.Each(
		func(o sku.Sku2) (err error) {
			errors.Out().Print(o)
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Last) runWithTransaktion(u *umwelt.Umwelt) (err error) {
	s := u.StoreObjekten()

	var transaktion *transaktion.Transaktion

	if transaktion, err = s.GetTransaktionStore().ReadLastTransaktion(); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.TodoP3("support log line format for skus")
	if err = transaktion.Skus.Each(
		func(o sku.SkuLike) (err error) {
			errors.Out().Print(o)
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
