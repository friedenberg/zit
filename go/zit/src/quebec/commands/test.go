package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type Test struct{}

func init() {
	registerCommandWithoutRepo(
		"test",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &Test{}

			return c
		},
	)
}

func (c Test) Run(
	u *repo_local.Repo,
	args ...string,
) (err error) {
	return
}
