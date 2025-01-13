package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local_working_copy"
)

type DormantRemove struct{}

func init() {
	registerCommand(
		"dormant-remove",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &DormantRemove{}

			return c
		},
	)
}

func (c DormantRemove) RunWithRepo(u *repo_local_working_copy.Repo, args ...string) {
	u.Must(u.Lock)

	for _, v := range args {
		cs := catgut.MakeFromString(v)

		if err := u.GetDormantIndex().RemoveDormantTag(cs); err != nil {
			u.CancelWithError(err)
		}
	}

	u.Must(u.Unlock)
}
