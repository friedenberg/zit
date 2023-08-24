package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/angeboren"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/papa/remote_transfers"
)

type Clone struct {
	Einleitung umwelt.Einleitung
}

func init() {
	registerCommandSansUmwelt(
		"clone",
		func(f *flag.FlagSet) Command {
			c := &Clone{
				Einleitung: umwelt.Einleitung{
					Angeboren: angeboren.Default(),
				},
			}

			c.Einleitung.AddToFlagSet(f)

			return c
		},
	)
}

func (c Clone) CompletionGattung() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Etikett,
		gattung.Zettel,
		gattung.Typ,
		gattung.Kasten,
	)
}

func (c Clone) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) == 0 {
		err = errors.Normalf("must specify kasten to pull from")
		return
	}

	from := args[0]

	if len(args) > 1 {
		args = args[1:]
	} else {
		err = errors.Normalf("Nothing to clone.")
		return
	}

	if err = u.Einleitung(c.Einleitung); err != nil {
		err = errors.Wrap(err)
		return
	}

	ids := u.MakeMetaIdSetWithExcludedHidden(
		matcher.MakeMatcherCwdNop(matcher.MakeMatcherAlways()),
		c.CompletionGattung(),
	)

	if err = ids.SetMany(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	var client remote_transfers.PullClient

	if client, err = remote_transfers.MakePullClient(u, from); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, client)

	if err = client.PullSkus(ids); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
