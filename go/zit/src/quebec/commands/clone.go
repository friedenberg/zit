package commands

import (
	"flag"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/bravo/xdg"
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
					Config:               immutable_config.Default(),
					ExcludeDefaultType:   true,
					ExcludeDefaultConfig: true,
				},
			}

			f.StringVar(&c.TheirXDGDotenv, "xdg-dotenv", "", "")

			c.BigBang.AddToFlagSet(f)

			return c
		},
	)
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

	var remote *env.Local

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

		if remote, err = c.cloneXDG(local, *dotenv.XDG); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		err = todo.Implement()
		return
	}

	ui.Debug().Print(remote)

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

	var list *sku.List

	if list, err = remote.MakeInventoryList(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	ui.Debug().Print(list)

	// get their inventory list as per the query in the args
	// setup the import to copy blobs from their env
	// import their inventory list

	return
}

func (c Clone) cloneXDG(
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
