package commands

import (
	"flag"
	"net"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Serve struct {
	Privileges string
}

func init() {
	registerCommand(
		"serve",
		func(f *flag.FlagSet) CommandWithContext {
			c := &Serve{}

			return c
		},
	)
}

func (c Serve) Run(u *env.Local, args ...string) {
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

	var listener net.Listener

	{
		var err error

		if listener, err = u.InitializeListener(network, address); err != nil {
			u.Context.Cancel(errors.Wrap(err))
			return
		}
	}

	defer u.Context.Closer(listener)

	if err := u.Serve(listener); err != nil {
		u.Context.Cancel(errors.Wrap(err))
		return
	}
}
