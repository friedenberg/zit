package commands

import (
	"encoding/json"
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

type Log struct {
}

func init() {
	registerCommand(
		"log",
		func(f *flag.FlagSet) Command {
			c := &Log{}

			return commandWithHinweisen{c}
		},
	)
}

func (c Log) RunWithHinweisen(os *umwelt.Umwelt, hs ...hinweis.Hinweis) (err error) {
	switch len(hs) {

	case 0:
		err = errors.Errorf("hinweis or zettel sha required")
		return
	}

	chains := make([]zettel_transacted.Slice, 0, len(hs))

	for _, h := range hs {
		var chain zettel_transacted.Slice

		if chain, err = os.StoreObjekten().AllInChain(h); err != nil {
			err = errors.Wrap(err)
			return
		}

		chains = append(chains, chain)
	}

	var b []byte

	if b, err = json.Marshal(chains); err != nil {
		err = errors.Wrapf(err, "failed to marshal json")
		return
	}

	errors.PrintOut(string(b))

	return
}
