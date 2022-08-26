package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/india/store_with_lock"
)

type Reindex struct {
}

func init() {
	registerCommand(
		"reindex",
		func(f *flag.FlagSet) Command {
			c := &Reindex{}

			return commandWithLockedStore{c}
		},
	)
}

func (c Reindex) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	if len(args) > 0 {
		err = errors.Errorf("reindex does not support arguments")
		return
	}

	logz.Print()

	if err = store.StoreObjekten().Reindex(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
