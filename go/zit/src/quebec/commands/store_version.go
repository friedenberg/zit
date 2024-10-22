package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type StoreVersion struct {
	immutable_config.Config
}

func init() {
	registerCommandWithoutEnvironment(
		"store-version",
		func(f *flag.FlagSet) CommandWithResult {
			c := &StoreVersion{
				Config: immutable_config.Default(),
			}

			return c
		},
	)
}

func (c StoreVersion) Run(u *env.Env, args ...string) (result Result) {
	ui.Out().Print(c.Config.GetStoreVersion())

	return
}
