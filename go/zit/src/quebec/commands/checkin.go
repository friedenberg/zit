package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Checkin struct {
	Delete     bool
	IgnoreBlob bool
	Organize   bool
}

func init() {
	registerCommandWithQuery(
		"checkin",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Checkin{}

			f.BoolVar(&c.Delete, "delete", false, "the checked-out file")

			f.BoolVar(
				&c.IgnoreBlob,
				"ignore-blob",
				false,
				"do not change the blob",
			)

			f.BoolVar(&c.Organize, "organize", false, "")

			return c
		},
	)
}

func (c Checkin) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.TrueGenre()...)
}

func (c Checkin) RunWithQuery(
	u *env.Env,
	qg *query.Group,
) (err error) {
	op := user_ops.Checkin{
		Delete: c.Delete,
	}

	if err = op.Run(
		u,
		qg,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
