package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/angeboren"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type StoreVersion struct {
	Angeboren angeboren.Konfig
}

func init() {
	registerCommandSansUmwelt(
		"store-version",
		func(f *flag.FlagSet) Command {
			c := &StoreVersion{
				Angeboren: angeboren.Default(),
			}

			return c
		},
	)
}

func (c StoreVersion) Run(u *umwelt.Umwelt, args ...string) (err error) {
	errors.Out().Printf("%s\n", c.Angeboren.GetStoreVersion())

	return
}
