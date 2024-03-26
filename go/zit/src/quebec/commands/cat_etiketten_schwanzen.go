package commands

import (
	"flag"
	"syscall"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type CatEtikettenSchwanzen struct{}

func init() {
	registerCommand(
		"cat-etiketten-schwanzen",
		func(f *flag.FlagSet) Command {
			c := &CatEtikettenSchwanzen{}

			return c
		},
	)
}

func (c CatEtikettenSchwanzen) Run(
	u *umwelt.Umwelt,
	_ ...string,
) (err error) {
	if err = u.GetStore().GetKennungIndex().EachSchwanzen(
		func(e *kennung.IndexedLike) (err error) {
			if err = errors.Out().Print(e.GetKennung().String()); err != nil {
				err = errors.IsAsNilOrWrapf(
					err,
					syscall.EPIPE,
					"Etikett: %s",
					e.GetKennung(),
				)

				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
