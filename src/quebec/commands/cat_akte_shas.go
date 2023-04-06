package commands

import (
	"flag"
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type CatAkteShas struct{}

func init() {
	registerCommand(
		"cat-akte-shas",
		func(f *flag.FlagSet) Command {
			c := &CatAkteShas{}

			return commandWithIds{c}
		},
	)
}

func (c CatAkteShas) CompletionGattung() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Akte,
	)
}

func (c CatAkteShas) RunWithIds(u *umwelt.Umwelt, ids kennung.Set) (err error) {
	if err = u.Standort().ReadAllShasForGattung(
		gattung.Akte,
		func(s sha.Sha) (err error) {
			_, err = fmt.Fprintln(u.Out(), s)
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
