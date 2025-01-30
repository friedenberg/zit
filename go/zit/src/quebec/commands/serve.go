package commands

import (
	"flag"
	"net"

	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/oscar/remote_http"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register("serve", &Serve{})
}

type Serve struct {
	command_components.Env
	command_components.EnvRepo
	command_components.LocalArchive
}

func (cmd *Serve) SetFlagSet(f *flag.FlagSet) {
	cmd.EnvRepo.SetFlagSet(f)
	cmd.LocalArchive.SetFlagSet(f)
}

func (cmd Serve) Run(req command.Request) {
	args := req.Args()
	req.SetCancelOnSIGHUP()

	envLocal := cmd.MakeEnvWithOptions(
		req,
		env_ui.Options{
			UIFileIsStderr: true,
			IgnoreTtyState: true,
		},
	)

	envRepo := cmd.MakeEnvRepoFromEnvLocal(envLocal)

	repo := cmd.MakeLocalArchive(envRepo)

	server := remote_http.Server{
		EnvLocal: envLocal,
		Repo:     repo,
	}

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
		server.ServeStdio()
	} else {
		var listener net.Listener

		{
			var err error

			if listener, err = server.InitializeListener(
				network,
				address,
			); err != nil {
				envLocal.CancelWithError(err)
			}

			defer envLocal.MustClose(listener)
		}

		if err := server.Serve(listener); err != nil {
			envLocal.CancelWithError(err)
		}
	}
}
