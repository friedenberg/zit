package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/alfa/angeboren"
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type Init struct {
	Einleitung umwelt.Einleitung
}

func init() {
	registerCommandSansUmwelt(
		"init",
		func(f *flag.FlagSet) Command {
			c := &Init{
				Einleitung: umwelt.Einleitung{
					Angeboren: angeboren.Default(),
				},
			}

			c.Einleitung.AddToFlagSet(f)

			return c
		},
	)
}

func (c Init) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if err = u.Einleitung(c.Einleitung); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
