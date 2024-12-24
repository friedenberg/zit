package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type DormantRemove struct{}

func init() {
	registerCommand(
		"schlummernd-remove",
		func(f *flag.FlagSet) Command {
			c := &DormantRemove{}

			return c
		},
	)
}

func (c DormantRemove) Run(u *env.Local, args ...string) (err error) {
	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, v := range args {
		cs := catgut.MakeFromString(v)

		if err = u.GetDormantIndex().RemoveDormantTag(cs); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	defer errors.Deferred(&err, u.Unlock)

	return
}
