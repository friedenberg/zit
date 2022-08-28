package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/charlie/zk_types"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/india/store_with_lock"
)

type Revert struct {
	Type zk_types.Type
}

func init() {
	registerCommand(
		"revert",
		func(f *flag.FlagSet) Command {
			c := &Revert{
				Type: zk_types.TypeUnknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithLockedStore{c}
		},
	)
}

func (c Revert) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	switch c.Type {
	case zk_types.TypeZettel:
		hins := make([]hinweis.Hinweis, len(args))

		for i, arg := range args {
			if hins[i], err = hinweis.Make(arg); err != nil {
				err = errors.Error(err)
				return
			}

			if _, err = store.StoreObjekten().Revert(hins[i]); err != nil {
				err = errors.Error(err)
				return
			}
		}

	default:
		stdprinter.Errf("objekte type %s does not support reverts currently\n", c.Type)
	}

	return
}
