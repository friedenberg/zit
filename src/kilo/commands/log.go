package commands

import (
	"encoding/json"
	"flag"

	"github.com/friedenberg/zit/collections"
	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/india/store_with_lock"
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
		stdprinter.Errf("ignoring extra arguments: %q\n", hs[1:])

		fallthrough

	case 1:
		h = hs[0]
	}

	var chain collections.SliceTransacted
	logz.Print()

	if chain, err = os.Zettels().AllInChain(h); err != nil {
		err = errors.Error(err)
		return
	}

	var b []byte

	if b, err = json.Marshal(chain); err != nil {
		err = errors.Wrapped(err, "failed to marshal json")
		return
	}

	stdprinter.Out(string(b))

	return
}
