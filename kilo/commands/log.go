package commands

import (
	"encoding/json"
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/hotel/zettels"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type Log struct {
}

func init() {
	registerCommand(
		"log",
		func(f *flag.FlagSet) Command {
			c := &Log{}

			return commandWithLockedStore{c}
		},
	)
}

func (c Log) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	var rawId string

	switch len(args) {

	case 0:
		err = errors.Errorf("hinweis or zettel sha required")
		return

	default:
		stdprinter.Errf("ignoring extra arguments: %q\n", args[1:])

		fallthrough

	case 1:
		rawId = args[0]

	}

	var id id.Id

	if id, err = c.getIdFromArg(rawId); err != nil {
		err = errors.Error(err)
		return
	}

	var chain zettels.Chain

	if chain, err = store.Zettels().AllInChain(id); err != nil {
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

func (c Log) getIdFromArg(arg string) (id id.Id, err error) {
	var sha sha.Sha

	if err = sha.Set(arg); err == nil {
		id = sha
		return
	}

	hinweis := hinweis.NewEmpty()

	if err = hinweis.Set(arg); err == nil {
		id = hinweis
		return
	}

	err = errors.Errorf("incorrect format for id: '%s'", arg)

	return
}
