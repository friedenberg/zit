package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Reindex struct {
	UseBestandsaufnahme bool
}

func init() {
	registerCommand(
		"reindex",
		func(f *flag.FlagSet) Command {
			c := &Reindex{}

			f.BoolVar(&c.UseBestandsaufnahme, "use-bestandsaufnahme", false, "use bestandsaufnahme")

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

	defer errors.Deferred(&err, u.Unlock)

	if c.UseBestandsaufnahme {
		if err = u.StoreObjekten().ReindexBestandsaufnahme(); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = u.StoreObjekten().Reindex(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
