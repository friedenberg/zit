package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Test struct{}

func init() {
	registerCommandWithoutEnvironment(
		"test",
		func(f *flag.FlagSet) Command {
			c := &Test{}

			return c
		},
	)
}

func (c Test) Run(
	u *env.Local,
	args ...string,
) (err error) {
	return
}
