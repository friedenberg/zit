package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/november/umwelt"
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

func (c Reindex) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) > 0 {
		err = errors.Errorf("reindex does not support arguments")
		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer u.Unlock()

	if err = u.StoreObjekten().Reindex(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}