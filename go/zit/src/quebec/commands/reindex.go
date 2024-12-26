package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type Reindex struct{}

func init() {
	registerCommand(
		"reindex",
		func(_ *flag.FlagSet) CommandWithRepo {
			c := &Reindex{}

			return c
		},
	)
}

func (c Reindex) GetEnvInitializeOptions() repo_local.Options {
	return repo_local.OptionsAllowConfigReadError
}

func (c Reindex) RunWithRepo(u *repo_local.Repo, args ...string) {
	if len(args) > 0 {
		u.CancelWithErrorf("reindex does not support arguments")
	}

	u.Must(u.Lock)
	u.Must(u.GetConfig().Reset)
	u.Must(u.GetStore().Reindex)
	u.Must(u.Unlock)
}
