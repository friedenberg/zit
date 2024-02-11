package commands

import (
	"flag"

	"code.linenisgreat.com/zit-go/src/alfa/angeboren"
	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/oscar/umwelt"
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
