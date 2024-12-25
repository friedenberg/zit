package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
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

func (c Reindex) GetEnvInitializeOptions() repo_local.Options {
	return repo_local.OptionsAllowConfigReadError
}

func (c Reindex) Run(u *repo_local.Local, args ...string) (err error) {
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
