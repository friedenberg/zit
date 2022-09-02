package commands

import (
	"encoding/json"
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/lima/store_with_lock"
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

func (c Log) RunWithHinweisen(os store_with_lock.Store, hs ...hinweis.Hinweis) (err error) {
	var h hinweis.Hinweis

	switch len(hs) {

	case 0:
		err = errors.Errorf("hinweis or zettel sha required")
		return

	default:
		errors.PrintErrf("ignoring extra arguments: %q", hs[1:])

		fallthrough

	case 1:
		h = hs[0]
	}

	var chain zettel_transacted.Slice
	errors.Print()

	if chain, err = os.StoreObjekten().AllInChain(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	var b []byte

	if b, err = json.Marshal(chain); err != nil {
		err = errors.Wrapf(err, "failed to marshal json")
		return
	}

	errors.PrintOut(string(b))

	return
}
