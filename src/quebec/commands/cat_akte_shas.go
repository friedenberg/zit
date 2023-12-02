package commands

import (
	"flag"
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type CatAkteShas struct{}

func init() {
	registerCommand(
		"cat-akte-shas",
		func(f *flag.FlagSet) Command {
			c := &CatAkteShas{}

			return c
		},
	)
}

func (c CatAkteShas) CompletionGattung() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Akte,
	)
}

func (c CatAkteShas) Run(u *umwelt.Umwelt, _ ...string) (err error) {
	if err = u.Standort().ReadAllShasForGattung(
		u.Konfig().GetStoreVersion(),
		gattung.Akte,
		func(s *sha.Sha) (err error) {
			_, err = fmt.Fprintln(u.Out(), s)
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
