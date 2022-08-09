package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type Rewrite struct {
}

func init() {
	registerCommand(
		"rewrite",
		func(f *flag.FlagSet) Command {
			c := &Rewrite{}

			return commandWithLockedStore{c}
		},
	)
}

func (c Rewrite) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	if len(args) > 0 {
		err = errors.Errorf("reindex does not support arguments")
		return
	}

	if err = store.Zettels().Rewrite(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
