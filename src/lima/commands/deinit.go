package commands

import (
	"flag"
	"os"
	"path"

	"github.com/friedenberg/zit/src/juliett/umwelt"
)

type Deinit struct {
  Force bool
}

func init() {
	registerCommand(
		"deinit",
		func(f *flag.FlagSet) Command {
			c := &Deinit{}

			return c
		},
	)
}

func (c Deinit) Run(u *umwelt.Umwelt, args ...string) (err error) {
  if !c.Force {
    return
  }

	base := path.Join(u.Standort().Dir(), ".zit")
	err = os.RemoveAll(base)

	if err != nil {
		return
	}

	return
}
