package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Copy struct {
	Edit bool
}

func init() {
	registerCommand(
		"cp",
		func(f *flag.FlagSet) Command {
			c := &Copy{}

			return commandWithIds{CommandWithIds: c}
		},
	)
}

func (c Copy) RunWithIds(s *umwelt.Umwelt, ids kennung.Set) (err error) {
	hins := ids.Hinweisen.Copy()

	zettels := make([]*zettel.Transacted, hins.Len())

	for i, h := range hins.Elements() {
		var tz *zettel.Transacted

		if tz, err = s.StoreObjekten().Zettel().ReadOne(h); err != nil {
			err = errors.Wrap(err)
			return
		}

		zettels[i] = tz
	}

	return
}
