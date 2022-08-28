package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/echo/umwelt"
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
