package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type StoreVersion struct {
	Angeboren immutable_config.Config
}

func init() {
	registerCommandSansUmwelt(
		"store-version",
		func(f *flag.FlagSet) Command {
			c := &StoreVersion{
				Angeboren: immutable_config.Default(),
			}

			return c
		},
	)
}

func (c StoreVersion) Run(u *umwelt.Umwelt, args ...string) (err error) {
	ui.Out().Printf("%s\n", c.Angeboren.GetStoreVersion())

	return
}
