package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type Reindex struct{}

func init() {
	registerCommand(
		"reindex",
		func(_ *flag.FlagSet) WithLocalWorkingCopy {
			c := &Reindex{}

			return c
		},
	)
}

func (c Reindex) GetLocalRepoOptions() local_working_copy.Options {
	return local_working_copy.OptionsAllowConfigReadError
}

func (c Reindex) RunWithLocalWorkingCopy(u *local_working_copy.Repo, args ...string) {
	if len(args) > 0 {
		u.CancelWithErrorf("reindex does not support arguments")
	}

	u.Must(u.Lock)
	u.Must(u.GetConfig().Reset)
	u.Must(u.GetStore().Reindex)
	u.Must(u.Unlock)
}
