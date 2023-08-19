package commands

import (
	"flag"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/november/umwelt"
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
	if err = u.StoreObjekten().GetKennungIndex().EachSchwanzen(
		func(e kennung.IndexedLike[kennung.Etikett, *kennung.Etikett]) (err error) {
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
