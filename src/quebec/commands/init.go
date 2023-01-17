package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/angeboren"
	"github.com/friedenberg/zit/src/oscar/umwelt"
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

			f.BoolVar(&c.Einleitung.DisableAge, "disable-age", false, "")
			f.StringVar(&c.Einleitung.Yin, "yin", "", "File containing list of Kennung")
			f.StringVar(&c.Einleitung.Yang, "yang", "", "File containing list of Kennung")
			c.Einleitung.Angeboren.AddToFlags(f)

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
