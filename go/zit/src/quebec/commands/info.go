package commands

import (
	"flag"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/xdg"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
)

type Info struct {
	immutable_config.Config
}

func init() {
	registerCommand(
		"info",
		func(f *flag.FlagSet) CommandWithEnv {
			c := &Info{
				Config: immutable_config.Default(),
			}

			return c
		},
	)
}

// TODO disambiguate this from repo / env
func (c Info) RunWithEnv(e *env.Env, args ...string) {
	if len(args) == 0 {
		args = []string{"store-version"}
	}

	for _, arg := range args {
		switch strings.ToLower(arg) {
		case "store-version":
			e.GetUI().Print(c.Config.GetStoreVersion())

		case "compression-type":
			e.GetUI().Print(c.Config.GetBlobStoreImmutableConfig().GetCompressionType())

		case "xdg":
			ecksDeeGee := e.GetDirLayout().GetXDG()

			dotenv := xdg.Dotenv{
				XDG: &ecksDeeGee,
			}

			if _, err := dotenv.WriteTo(e.GetOutFile()); err != nil {
				e.CancelWithError(err)
			}
		}
	}
}
