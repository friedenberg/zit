package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/november/repo_local_working_copy"
)

type Test struct{}

func init() {
	registerCommand(
		"test",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &Test{}

			return c
		},
	)
}

func (c Test) RunWithRepo(
	u *repo_local_working_copy.Repo,
	args ...string,
) {
}
