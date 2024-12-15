package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Clone struct {
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

			c.BigBang.AddToFlagSet(f)

			return c
		},
	)
}

func (c Clone) Run(u *env.Env, args ...string) (err error) {
	if err = u.Start(c.BigBang); err != nil {
		err = errors.Wrap(err)
		return
	}

	if len(args) < 1 {
		// TODO add info about remote options
		err = errors.BadRequestf("Cloning requires a remote to be specified")
		return
	}

	return
}
