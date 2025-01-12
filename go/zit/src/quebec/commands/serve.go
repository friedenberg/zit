package commands

import (
	"flag"
	"net"

	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/read_write_repo_local"
)

type Serve struct {
	Privileges string
}

func init() {
	registerCommand(
		"serve",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &Serve{}

			return c
		},
	)
}

func (c Serve) GetEnvOptions() env.Options {
	return env.Options{
		UIFileIsStderr: true,
	}
}

func (c Serve) RunWithRepo(u *read_write_repo_local.Repo, args ...string) {
	u.SetCancelOnSIGHUP()

	// TODO switch network to be RemoteServeType
	var network, address string

	switch len(args) {
	case 0:
		network = "tcp"
		address = ":0"

	case 1:
		network = args[0]

	default:
		network = args[0]
		address = args[1]
	}

	if network == "-" {
		if err := u.ServeStdio(); err != nil {
			u.CancelWithError(err)
		}
	} else {
		var listener net.Listener

		{
			var err error

			if listener, err = u.InitializeListener(network, address); err != nil {
				u.CancelWithError(err)
			}

			defer u.MustClose(listener)
		}

		if err := u.Serve(listener); err != nil {
			u.CancelWithError(err)
		}
	}
}
