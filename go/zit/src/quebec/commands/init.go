package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type Init struct {
	repo_local.BigBang
}

func init() {
	registerCommandWithoutRepo(
		"init",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &Init{
				BigBang: repo_local.BigBang{
					Config: immutable_config.Default(),
				},
			}

			c.BigBang.AddToFlagSet(f)

			return c
		},
	)
}

func (c Init) RunWithRepo(u *repo_local.Repo, args ...string) {
	if err := u.Start(c.BigBang); err != nil {
		u.CancelWithError(err)
		return
	}

	return
}
