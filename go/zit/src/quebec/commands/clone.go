package commands

import (
	"flag"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/bravo/xdg"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
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

		if remote, err = c.cloneXDG(local, *dotenv.XDG, args...); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		err = todo.Implement()
		return
	}

	ui.Debug().Print(remote)
	// get their inventory list as per the query in the args
	// setup the import to copy blobs from their env
	// import their inventory list

	return
}

func (c Clone) cloneXDG(
	local *env.Local,
	ecksDeeGee xdg.XDG,
	args ...string,
) (remote *env.Local, err error) {
	// if remote, err = env.
	// bootstrap their dirlayout and turn it into an *env.Env

	return
}
