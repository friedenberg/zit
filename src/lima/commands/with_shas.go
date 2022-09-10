package commands

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/juliett/umwelt"
)

type WithShas interface {
	RunWithShas(store *umwelt.Umwelt, shas ...sha.Sha) error
}

type withShas struct {
	WithShas
}

func (c withShas) Run(store *umwelt.Umwelt, args ...string) (err error) {
	shas := make([]sha.Sha, len(args))

	for i, arg := range args {
		var sha sha.Sha

		if err = sha.Set(arg); err != nil {
			err = errors.Wrap(err)
			return
		}

		shas[i] = sha
	}

	if err = c.RunWithShas(store, shas...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
