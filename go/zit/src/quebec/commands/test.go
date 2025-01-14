package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type Test struct{}

func init() {
	registerCommand(
		"test",
		func(f *flag.FlagSet) WithLocalWorkingCopy {
			c := &Test{}

			return c
		},
	)
}

func (c Test) Run(
	u *local_working_copy.Repo,
	args ...string,
) {
}
