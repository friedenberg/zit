package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

// TODO-P1 determine if this can be removed and replaced with show
type ExpandZettelId struct{}

func init() {
	registerCommand(
		"expand-hinweis",
		func(f *flag.FlagSet) Command {
			c := &ExpandZettelId{}

			return c
		},
	)
}

func (c ExpandZettelId) Run(u *env.Env, args ...string) (err error) {
	for _, v := range args {
		var h *ids.ZettelId

		h, err = u.GetStore().GetAbbrStore().ZettelId().ExpandString(v)
		if err != nil {
			err = errors.Wrap(err)
			return
		}

		ui.Out().Print(h)
	}

	return
}
