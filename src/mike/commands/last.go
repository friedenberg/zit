package commands

import (
	"encoding/json"
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/zk_types"
	"github.com/friedenberg/zit/src/delta/transaktion"
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

type Last struct {
	Type zk_types.Type
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

	errors.PrintDebug(transaktion.Objekten)

	enc := json.NewEncoder(u.Out())

	if err = enc.Encode(transaktion); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
