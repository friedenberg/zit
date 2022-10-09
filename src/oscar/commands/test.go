package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/mike/umwelt"
)

type Test struct {
}

func init() {
	registerCommand(
		"test",
		func(f *flag.FlagSet) Command {
			c := &Test{}

			return c
		},
	)
}

func (c Test) Run(u *umwelt.Umwelt, args ...string) (err error) {
	return
}
