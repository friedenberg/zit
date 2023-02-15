package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type ExpandHinweis struct{}

func init() {
	registerCommand(
		"expand-hinweis",
		func(f *flag.FlagSet) Command {
			c := &ExpandHinweis{}

			return commandWithIds{CommandWithIds: c}
		},
	)
}

func (c ExpandHinweis) CompletionGattung() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Zettel,
	)
}

func (c ExpandHinweis) RunWithIds(s *umwelt.Umwelt, ids kennung.Set) (err error) {
	hins := ids.Hinweisen.ImmutableClone()

	for _, h := range hins.Elements() {
		errors.Out().Print(h)
	}

	return
}
