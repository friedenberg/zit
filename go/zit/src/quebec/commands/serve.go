package commands

import (
	"net"

	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	registerCommand("serve", &Serve{})
}

type Serve struct {
	command_components.LocalWorkingCopy
}

func (c Serve) Run(dep command.Dep) {
	args := dep.Args()
	dep.SetCancelOnSIGHUP()

	localWorkingCopy := c.MakeLocalWorkingCopyWithOptions(
		dep,
		env.Options{
			UIFileIsStderr: true,
		},
		local_working_copy.OptionsEmpty,
	)

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
		if err := localWorkingCopy.ServeStdio(); err != nil {
			localWorkingCopy.CancelWithError(err)
		}
	} else {
		var listener net.Listener

		{
			var err error

			if listener, err = localWorkingCopy.InitializeListener(network, address); err != nil {
				localWorkingCopy.CancelWithError(err)
			}

			defer localWorkingCopy.MustClose(listener)
		}

		if err := localWorkingCopy.Serve(listener); err != nil {
			localWorkingCopy.CancelWithError(err)
		}
	}
}
