package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type DormantAdd struct{}

func init() {
	registerCommand(
		"dormant-add",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &DormantAdd{}

			return c
		},
	)
}

func (c DormantAdd) RunWithRepo(u *repo_local.Repo, args ...string) {
	u.Must(u.Lock)

	for _, v := range args {
		cs := catgut.MakeFromString(v)

		if err := u.GetDormantIndex().AddDormantTag(cs); err != nil {
			u.CancelWithError(err)
		}
	}

	u.Must(u.Unlock)
}
