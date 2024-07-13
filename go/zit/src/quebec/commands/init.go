package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
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
					Angeboren: immutable_config.Default(),
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
