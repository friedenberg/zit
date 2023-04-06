package commands

import (
	"flag"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type CatEtikettenSchwanzen struct{}

func init() {
	registerCommand(
		"cat-etiketten-schwanzen",
		func(f *flag.FlagSet) Command {
			c := &CatEtikettenSchwanzen{}

			return commandWithIds{c}
		},
	)
}

func (c CatEtikettenSchwanzen) CompletionGattung() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Etikett,
	)
}

func (c CatEtikettenSchwanzen) RunWithIds(u *umwelt.Umwelt, ids kennung.Set) (err error) {
	var ea []kennung.Etikett

	if ea, err = u.StoreObjekten().GetKennungIndex().GetAllEtiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, e := range ea {
		if err = errors.Out().Print(e.String()); err != nil {
			err = errors.IsAsNilOrWrapf(
				err,
				syscall.EPIPE,
				"Etikett: %s",
				e,
			)

			return
		}
	}

	return
}
