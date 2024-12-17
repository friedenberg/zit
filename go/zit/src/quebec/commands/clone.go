package commands

import (
	"flag"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/xdg"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Clone struct {
	TheirXDGDotenv string
	env.BigBang
}

func init() {
	registerCommandWithoutEnvironment(
		"clone",
		func(f *flag.FlagSet) Command {
			c := &Clone{
				BigBang: env.BigBang{
					Config:             immutable_config.Default(),
					ExcludeDefaultType: true,
				},
			}

			f.StringVar(&c.TheirXDGDotenv, "xdg-dotenv", "", "")

			c.BigBang.AddToFlagSet(f)

			return c
		},
	)
}

func (c Clone) DefaultSigil() ids.Sigil {
	return ids.MakeSigil(ids.SigilHistory, ids.SigilHidden)
}

func (c Clone) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.InventoryList)
	// return ids.MakeGenre(genres.TrueGenre()...)
}

func (c Clone) Run(local *env.Local, args ...string) (err error) {
	if err = local.Start(c.BigBang); err != nil {
		err = errors.Wrap(err)
		return
	}

	if len(args) < 1 && c.TheirXDGDotenv == "" {
		// TODO add info about remote options
		err = errors.BadRequestf("Cloning requires a remote to be specified")
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

		if remote, err = c.makeLocalEnvFromXDG(local, *dotenv.XDG); err != nil {
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

	if err = local.PullQueryGroupFromRemote(
		remote,
		qg,
		true,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Clone) makeLocalEnvFromXDG(
	local *env.Local,
	xdg xdg.XDG,
) (remote *env.Local, err error) {
	var primitiveFSHome dir_layout.Primitive

	if primitiveFSHome, err = dir_layout.MakePrimitiveWithXDG(
		local.GetConfig().Debug,
		xdg,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if remote, err = env.Make(
		nil,
		local.GetConfig().Cli(),
		env.OptionsEmpty,
		primitiveFSHome,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
