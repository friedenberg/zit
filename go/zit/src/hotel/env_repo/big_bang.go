package env_repo

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
)

type BigBang struct {
	ids.Type
	Config *immutable_config.Latest

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
	bb.Config = immutable_config.Default()
	bb.Config.SetFlagSet(f)
}
