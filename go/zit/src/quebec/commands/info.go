package commands

import (
	"flag"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/delta/config_immutable"
	"code.linenisgreat.com/zit/go/zit/src/delta/xdg"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
)

type Info struct {
	config_immutable.Config
}

func init() {
	command.Register(
		"info",
		&Info{
			Config: config_immutable.Default(),
		},
	)
}

func (c Info) SetFlagSet(f *flag.FlagSet) {}

func (c Info) Run(dependencies command.Request) {
	dir := env_dir.MakeDefault(
		dependencies,
		dependencies.Debug,
	)

	ui := env_ui.Make(
		dependencies,
		dependencies.Config,
		env_ui.Options{},
	)

	args := dependencies.Args()

	if len(args) == 0 {
		args = []string{"store-version"}
	}

	for _, arg := range args {
		switch strings.ToLower(arg) {
		case "store-version":
			ui.GetUI().Print(c.Config.GetStoreVersion())

		case "compression-type":
			ui.GetUI().Print(c.Config.GetBlobStoreImmutableConfig().GetCompressionType())

		case "age-encryption":
			ui.GetUI().Print(c.Config.GetBlobStoreImmutableConfig().GetAgeEncryption())

		case "xdg":
			ecksDeeGee := dir.GetXDG()

			dotenv := xdg.Dotenv{
				XDG: &ecksDeeGee,
			}

			if _, err := dotenv.WriteTo(ui.GetOutFile()); err != nil {
				ui.CancelWithError(err)
			}
		}
	}
}
