package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/india/query"
	"code.linenisgreat.com/zit/src/kilo/cwd"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
	"code.linenisgreat.com/zit/src/papa/user_ops"
)

type Checkin struct {
	Delete     bool
	IgnoreAkte bool
}

func init() {
	registerCommandWithCwdQuery(
		"checkin",
		func(f *flag.FlagSet) CommandWithCwdQuery {
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

func (c Checkin) RunWithCwdQuery(
	u *umwelt.Umwelt,
	ms *query.Group,
	pz *cwd.CwdFiles,
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
