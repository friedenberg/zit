package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Reindex struct{}

func init() {
	registerCommand(
		"reindex",
		func(_ *flag.FlagSet) Command {
			c := &Reindex{}

			return c
		},
	)
}

func (c Reindex) GetUmweltInitializeOptions() umwelt.Options {
	return umwelt.OptionsAllowKonfigReadError
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

	defer errors.Deferred(&err, u.Unlock)

	if err = u.StoreObjekten().Reindex(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}