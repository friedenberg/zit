package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/juliett/umwelt"
)

type Reindex struct {
}

func init() {
	registerCommand(
		"reindex",
		func(f *flag.FlagSet) Command {
			c := &Reindex{}

			return c
		},
	)
}

func (c Reindex) Run(store *umwelt.Umwelt, args ...string) (err error) {
	if len(args) > 0 {
		err = errors.Errorf("reindex does not support arguments")
		return
	}

	errors.Print()

	if err = store.StoreObjekten().Reindex(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
