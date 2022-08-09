package commands

import (
	"encoding/json"
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/hotel/zettels"
)

type Log struct {
}

func init() {
	registerCommand(
		"log",
		func(f *flag.FlagSet) Command {
			c := &Log{}

			return commandWithLockedStore{commandWithHinweisen{c}}
		},
	)
}

func (c Log) RunWithHinweisen(u *umwelt.Umwelt, zs zettels.Zettels, hs ...hinweis.Hinweis) (err error) {
	var h hinweis.Hinweis

	switch len(hs) {

	case 0:
		err = errors.Errorf("hinweis or zettel sha required")
		return

	default:
		stdprinter.Errf("ignoring extra arguments: %q\n", hs[1:])

		fallthrough

	case 1:
		h = hs[0]
	}

	var chain zettels.Chain
	logz.Print()

	if chain, err = zs.AllInChain(h); err != nil {
		err = errors.Error(err)
		return
	}

	b, err := json.Marshal(chain)

	if err != nil {
		logz.Print(err)
	} else {
		stdprinter.Out(string(b))
	}

	return
}
