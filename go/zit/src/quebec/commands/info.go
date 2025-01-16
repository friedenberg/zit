package commands

import (
	"flag"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/xdg"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
)

type Info struct {
	immutable_config.Config
}

func init() {
	registerCommand(
		"info",
		&Info{
			Config: immutable_config.Default(),
		},
	)
}

func (c Info) SetFlagSet(f *flag.FlagSet) {}

func (c Info) Run(dependencies Dependencies) {
	layout := dir_layout.MakeDefault(
		dependencies.Context,
		dependencies.Debug,
	)

	env := env.Make(
		dependencies.Context,
		dependencies.Config,
		layout,
		env.Options{},
	)

	args := dependencies.Args()

	if len(args) == 0 {
		args = []string{"store-version"}
	}

	for _, arg := range args {
		switch strings.ToLower(arg) {
		case "store-version":
			env.GetUI().Print(c.Config.GetStoreVersion())

		case "compression-type":
			env.GetUI().Print(c.Config.GetBlobStoreImmutableConfig().GetCompressionType())

		case "age-encryption":
			env.GetUI().Print(c.Config.GetBlobStoreImmutableConfig().GetAgeEncryption())

		case "xdg":
			ecksDeeGee := env.GetDirLayout().GetXDG()

			dotenv := xdg.Dotenv{
				XDG: &ecksDeeGee,
			}

			if _, err := dotenv.WriteTo(env.GetOutFile()); err != nil {
				env.CancelWithError(err)
			}
		}
	}
}
