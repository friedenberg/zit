package commands

import (
	"flag"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/bravo/xdg"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Info struct {
	immutable_config.Config
}

func init() {
	registerCommandWithoutEnvironment(
		"info",
		func(f *flag.FlagSet) CommandWithResult {
			c := &Info{
				Config: immutable_config.Default(),
			}

			return c
		},
	)
}

func (c Info) Run(u *env.Local, args ...string) (result Result) {
	if len(args) == 0 {
		args = []string{"store-version"}
	}

	for _, arg := range args {
		switch strings.ToLower(arg) {
		case "store-version":
			ui.Out().Print(c.Config.GetStoreVersion())

		case "xdg":
			ecksDeeGee := u.GetDirLayoutPrimitive().GetXDG()

			dotenv := xdg.Dotenv{
				XDG: &ecksDeeGee,
			}

			if _, result.Error = dotenv.WriteTo(u.Out()); result.Error != nil {
				result.Error = errors.Wrap(result.Error)
				return
			}
		}
	}

	return
}
