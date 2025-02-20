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
	config_immutable.ConfigPrivate
}

func init() {
	command.Register(
		"info",
		&Info{
			ConfigPrivate: config_immutable.Default(),
		},
	)
}

func (c Info) SetFlagSet(f *flag.FlagSet) {}

func (c Info) Run(req command.Request) {
	dir := env_dir.MakeDefault(
		req,
		req.Debug,
	)

	ui := env_ui.Make(
		req,
		req.Config,
		env_ui.Options{},
	)

	args := req.PopArgs()

	if len(args) == 0 {
		args = []string{"store-version"}
	}

	for _, arg := range args {
		switch strings.ToLower(arg) {
		case "store-version":
			ui.GetUI().Print(c.ConfigPrivate.GetStoreVersion())

		case "compression-type":
			ui.GetUI().Print(c.ConfigPrivate.GetBlobStoreConfigImmutable().GetBlobCompression())

		case "age-encryption":
			ui.GetUI().Print(c.ConfigPrivate.GetBlobStoreConfigImmutable().GetBlobEncryption())

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
