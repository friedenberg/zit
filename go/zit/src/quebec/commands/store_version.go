package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/angeboren"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
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
	ui.Out().Printf("%s\n", c.Angeboren.GetStoreVersion())

	return
}
