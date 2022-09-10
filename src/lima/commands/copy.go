package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/umwelt"
)

type Copy struct {
	Edit bool
}

func init() {
	registerCommand(
		"cp",
		func(f *flag.FlagSet) Command {
			c := &Copy{}

			return commandWithHinweisen{c}
		},
	)
}

func (c Copy) RunWithHinweisen(s *umwelt.Umwelt, hins ...hinweis.Hinweis) (err error) {
	zettels := make([]zettel_transacted.Zettel, len(hins))

	for i, h := range hins {
		var tz zettel_transacted.Zettel

		if tz, err = s.StoreObjekten().Read(h); err != nil {
			err = errors.Wrap(err)
			return
		}

		zettels[i] = tz
	}

	return
}
