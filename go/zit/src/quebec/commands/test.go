package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type Test struct{}

func init() {
	registerCommand(
		"test",
		func(_ *flag.FlagSet) Command {
			c := &Test{}

			return c
		},
	)
}

func (c Test) Run(u *umwelt.Umwelt, args ...string) (err error) {
	return u.GetStore().GetVerzeichnisse().GetEnnuiShas().PrintAll()
}
