package commands

import (
	"flag"
	"os"
	"path"
)

type Deinit struct {
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

func (c Deinit) Run(u _Umwelt, args ...string) (err error) {
	u.Lock.Lock()
	defer _PanicIfError(u.Lock.Unlock())

	base := path.Join(u.Dir(), ".zit")
	err = os.RemoveAll(base)

	if err != nil {
		return
	}

	return
}
