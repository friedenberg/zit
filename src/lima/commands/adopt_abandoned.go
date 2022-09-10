package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/umwelt"
)

type AdoptAbandoned struct {
}

func init() {
	registerCommand(
		"adopt-abandoned",
		func(f *flag.FlagSet) Command {
			c := &AdoptAbandoned{}

			return c
		},
	)
}

func (c AdoptAbandoned) Description() string {
	return "creates a new hinweis for a zettel that has somehow gotten detached"
}

func (c AdoptAbandoned) Run(u *umwelt.Umwelt, args ...string) (err error) {
	zettels := make([]zettel_named.Zettel, len(args))

	for i, a := range args {
		var sha sha.Sha

		if err = sha.Set(a); err != nil {
			err = errors.Wrap(err)
			return
		}

		var stored zettel_transacted.Zettel

		if stored, err = u.StoreObjekten().Read(sha); err != nil {
			err = errors.Wrap(err)
			return
		}

		var tz zettel_transacted.Zettel

		if tz, err = u.StoreObjekten().Create(stored.Named.Stored.Zettel); err != nil {
			err = errors.Wrap(err)
			return
		}

		zettels[i] = tz.Named
		errors.PrintOutf("%s (adopted)", tz.Named)
	}

	return
}
