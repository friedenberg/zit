package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/juliett/query"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
	"code.linenisgreat.com/zit/src/papa/user_ops"
)

type Checkin struct {
	Delete     bool
	IgnoreAkte bool
}

func init() {
	registerCommandWithQuery(
		"checkin",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Checkin{}

			f.BoolVar(&c.Delete, "delete", false, "the checked-out file")
			f.BoolVar(
				&c.IgnoreAkte,
				"ignore-akte",
				false,
				"do not change the akte",
			)

			return c
		},
	)
}

func (c Checkin) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung()
}

func (c Checkin) RunWithQuery(
	u *umwelt.Umwelt,
	ms *query.Group,
) (err error) {
	op := user_ops.Checkin{
		Delete: c.Delete,
	}

	if err = op.Run(u, ms); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
