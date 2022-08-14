package commands

import (
	"encoding/json"
	"flag"

	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/stdprinter"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/india/store_with_lock"
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

	var chain []stored_zettel.Transacted
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
