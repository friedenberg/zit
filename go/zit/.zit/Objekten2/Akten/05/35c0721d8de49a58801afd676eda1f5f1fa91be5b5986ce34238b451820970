package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type DormantAdd struct{}

func init() {
	registerCommand(
		"schlummernd-add",
		func(f *flag.FlagSet) Command {
			c := &DormantAdd{}

			return c
		},
	)
}

func (c DormantAdd) Run(u *env.Local, args ...string) (err error) {
	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, v := range args {
		cs := catgut.MakeFromString(v)

		if err = u.GetDormantIndex().AddDormantTag(cs); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	defer errors.Deferred(&err, u.Unlock)

	return
}
