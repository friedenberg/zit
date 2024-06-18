package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type SchlummerndAdd struct{}

func init() {
	registerCommand(
		"schlummernd-add",
		func(f *flag.FlagSet) Command {
			c := &SchlummerndAdd{}

			return c
		},
	)
}

func (c SchlummerndAdd) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, v := range args {
		cs := catgut.MakeFromString(v)

		if err = u.Schlummernd().AddSchlummerndEtikett(cs); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	defer errors.Deferred(&err, u.Unlock)

	return
}
