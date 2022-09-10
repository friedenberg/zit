package commands

import (
	"flag"
	"sort"
	"strconv"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/juliett/umwelt"
)

type PeekHinweisen struct {
}

func init() {
	registerCommand(
		"peek-hinweisen",
		func(f *flag.FlagSet) Command {
			c := &PeekHinweisen{}

			return c
		},
	)
}

func (c PeekHinweisen) Run(store *umwelt.Umwelt, args ...string) (err error) {
	n := 0

	if len(args) > 0 {
		if n, err = strconv.Atoi(args[0]); err != nil {
			err = errors.Errorf("expected int but got %s", args[0])
			return
		}
	}

	var hs []hinweis.Hinweis

	if hs, err = store.StoreObjekten().PeekHinweisen(n); err != nil {
		err = errors.Wrap(err)
		return
	}

	sort.Slice(
		hs,
		func(i, j int) bool {
			return hs[i].String() < hs[j].String()
		},
	)

	for i, h := range hs {
		errors.PrintOutf("%d: %s", i, h)
	}

	return
}
