package commands

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/juliett/umwelt"
	"github.com/friedenberg/zit/src/kilo/user_ops"
)

type CommandWithHinweisen interface {
	RunWithHinweisen(*umwelt.Umwelt, ...hinweis.Hinweis) error
}

type commandWithHinweisen struct {
	CommandWithHinweisen
}

func (c commandWithHinweisen) Run(
	store *umwelt.Umwelt,
	args ...string,
) (err error) {
	var hins []hinweis.Hinweis

	if hins, err = (user_ops.GetHinweisenFromArgs{}).RunMany(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.RunWithHinweisen(store, hins...)

	return
}
