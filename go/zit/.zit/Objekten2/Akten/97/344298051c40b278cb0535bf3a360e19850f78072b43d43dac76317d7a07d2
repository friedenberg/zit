package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Reindex struct{}

func init() {
	registerCommand(
		"reindex",
		func(_ *flag.FlagSet) Command {
			c := &Reindex{}

			return c
		},
	)
}

func (c Reindex) GetEnvInitializeOptions() env.Options {
	return env.OptionsAllowConfigReadError
}

func (c Reindex) Run(u *env.Local, args ...string) (err error) {
	if len(args) > 0 {
		err = errors.Errorf("reindex does not support arguments")
		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.GetConfig().Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.GetStore().Reindex(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
