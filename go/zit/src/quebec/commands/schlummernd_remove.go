package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type DormantRemove struct{}

func init() {
	registerCommand(
		"schlummernd-remove",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &DormantRemove{}

			return c
		},
	)
}

func (c DormantRemove) RunWithRepo(u *repo_local.Repo, args ...string) {
	if err := u.Lock(); err != nil {
		u.CancelWithError(err)
		return
	}

	for _, v := range args {
		cs := catgut.MakeFromString(v)

		if err := u.GetDormantIndex().RemoveDormantTag(cs); err != nil {
			u.CancelWithError(err)
			return
		}
	}

	u.Context.Must(u.Unlock)
}
