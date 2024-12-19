package commands

import (
	"flag"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/xdg"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Push struct {
	TheirXDGDotenv string
}

func init() {
	registerCommand(
		"push",
		func(f *flag.FlagSet) Command {
			c := &Push{}

			f.StringVar(&c.TheirXDGDotenv, "xdg-dotenv", "", "")

			return c
		},
	)
}

func (c Push) DefaultSigil() ids.Sigil {
	return ids.MakeSigil(ids.SigilHistory, ids.SigilHidden)
}

func (c Push) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.InventoryList)
}

func (c Push) Run(local *env.Local, args ...string) (err error) {
	if len(args) < 1 && c.TheirXDGDotenv == "" {
		// TODO add info about remote options
		err = errors.BadRequestf("Pushing requires a remote to be specified")
		return
	}

	var remote env.Env

	if c.TheirXDGDotenv != "" {
		dotenv := xdg.Dotenv{
			XDG: &xdg.XDG{},
		}

		var f *os.File

		if f, err = os.Open(c.TheirXDGDotenv); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = dotenv.ReadFrom(f); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = f.Close(); err != nil {
			err = errors.Wrap(err)
			return
		}

		if remote, err = env.MakeLocalFromConfigAndXDG(
      local.Context,
			local.GetConfig(),
			*dotenv.XDG,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		err = todo.Implement()
		return
	}

	var qg *query.Group

	if qg, err = remote.MakeQueryGroup(
		c,
		ids.RepoId{},
		sku.ExternalQueryOptions{},
		args...,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = remote.PullQueryGroupFromRemote(
		local,
		qg,
		true,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
