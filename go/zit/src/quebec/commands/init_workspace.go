package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_mutable_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/workspace_config_blobs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register(
		"init-workspace",
		&InitWorkspace{},
	)
}

type InitWorkspace struct {
	command_components.LocalWorkingCopy
	DefaultQueryGroup string
	Proto             sku.Proto
}

func (cmd *InitWorkspace) SetFlagSet(f *flag.FlagSet) {
	cmd.LocalWorkingCopy.SetFlagSet(f)
	cmd.Proto.SetFlagSet(f)
	f.StringVar(&cmd.DefaultQueryGroup, "query", "", "")
}

func (cmd InitWorkspace) Run(req command.Request) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(req)

	blob := &workspace_config_blobs.V0{
		Query: cmd.DefaultQueryGroup,
		Defaults: config_mutable_blobs.DefaultsV1OmitEmpty{
			Type: cmd.Proto.Type,
			Tags: quiter.Elements(cmd.Proto.Tags),
		},
	}

	if err := localWorkingCopy.GetEnvWorkspace().CreateWorkspace(
		blob,
	); err != nil {
		req.CancelWithError(err)
	}
}
