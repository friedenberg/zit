package env_repo

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/config_immutable"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
)

type BigBang struct {
	ids.Type
	Config *config_immutable.LatestPrivate

	Yin                  string
	Yang                 string
	ExcludeDefaultType   bool
	ExcludeDefaultConfig bool
	OverrideXDGWithCwd   bool
}

func (bb *BigBang) SetFlagSet(f *flag.FlagSet) {
	f.BoolVar(&bb.OverrideXDGWithCwd, "override-xdg-with-cwd", false, "")
	f.StringVar(&bb.Yin, "yin", "", "File containing list of zettel id left parts")
	f.StringVar(&bb.Yang, "yang", "", "File containing list of zettel id right parts")

	bb.Type = builtin_types.GetOrPanic(builtin_types.ImmutableConfigV1).Type
	bb.Config = config_immutable.Default()
	bb.Config.SetFlagSet(f)
}
