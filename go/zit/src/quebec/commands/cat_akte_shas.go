package commands

import (
	"flag"
	"fmt"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/delta/gattungen"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
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
