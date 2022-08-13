package commands

import (
	"flag"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/stdprinter"
	"github.com/friedenberg/zit/bravo/node_type"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type Revert struct {
	Type node_type.Type
}

func init() {
	registerCommand(
		"revert",
		func(f *flag.FlagSet) Command {
			c := &Revert{
				Type: node_type.TypeUnknown,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithLockedStore{c}
		},
	)
}

func (c Revert) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	switch c.Type {
	case node_type.TypeZettel:
		hins := make([]hinweis.Hinweis, len(args))

		for i, arg := range args {
			if hins[i], err = hinweis.MakeBlindHinweis(arg); err != nil {
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
