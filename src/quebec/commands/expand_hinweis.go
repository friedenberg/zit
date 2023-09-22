package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/november/umwelt"
)

// TODO-P1 determine if this can be removed and replaced with show
type ExpandHinweis struct{}

func init() {
	registerCommand(
		"expand-hinweis",
		func(f *flag.FlagSet) Command {
			c := &ExpandHinweis{}

			return c
		},
	)
}

func (c ExpandHinweis) Run(u *umwelt.Umwelt, args ...string) (err error) {
	for _, v := range args {
		var h kennung.Hinweis

		h, err = u.StoreObjekten().GetAbbrStore().Hinweis().ExpandString(v)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		errors.Out().Print(h)
	}

	return
}
