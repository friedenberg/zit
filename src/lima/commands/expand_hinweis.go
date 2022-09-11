package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/juliett/umwelt"
)

type ExpandHinweis struct {
}

func init() {
	registerCommand(
		"expand-hinweis",
		func(f *flag.FlagSet) Command {
			c := &ExpandHinweis{}

			return commandWithHinweisen{c}
		},
	)
}

func (c ExpandHinweis) RunWithHinweisen(s *umwelt.Umwelt, hins ...hinweis.Hinweis) (err error) {
	for _, h := range hins {
		errors.PrintOut(h)
	}

	return
}
