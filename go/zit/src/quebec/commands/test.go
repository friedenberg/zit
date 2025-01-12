package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/november/read_write_repo_local"
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
	u *read_write_repo_local.Repo,
	args ...string,
) {
}
