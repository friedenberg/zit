package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type Init struct {
	repo_local.BigBang
}

func init() {
	registerCommandWithoutEnvironment(
		"init",
		func(f *flag.FlagSet) Command {
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

func (c Init) Run(u *repo_local.Repo, args ...string) (err error) {
	if err = u.Start(c.BigBang); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
