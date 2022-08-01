package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type Revert struct {
	Type _Type
}

func init() {
	registerCommand(
		"revert",
		func(f *flag.FlagSet) Command {
			c := &Revert{
				Type: _TypeUnknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithLockedStore{c}
		},
	)
}

func (c Revert) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	switch c.Type {
	case _TypeZettel:
		hins := make([]_Hinweis, len(args))

		for i, arg := range args {
			if hins[i], err = _MakeBlindHinweis(arg); err != nil {
				err = errors.Error(err)
				return
			}

			if _, err = store.Zettels().Revert(hins[i]); err != nil {
				err = errors.Error(err)
				return
			}
		}

	default:
		stdprinter.Errf("objekte type %s does not support reverts currently\n", c.Type)
	}

	return
}
