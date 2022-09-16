package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/zk_types"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/juliett/umwelt"
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

			return c
		},
	)
}

func (c Revert) Run(store *umwelt.Umwelt, args ...string) (err error) {
	switch c.Type {
	case zk_types.TypeZettel:
		hins := make([]hinweis.Hinweis, len(args))

		if err = store.Lock(); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer store.Unlock()

		for i, arg := range args {
			if hins[i], err = hinweis.Make(arg); err != nil {
				err = errors.Wrap(err)
				return
			}

			if _, err = store.StoreObjekten().Revert(hins[i]); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	default:
		errors.PrintErrf("objekte type %s does not support reverts currently", c.Type)
	}

	return
}
