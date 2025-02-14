package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register(
		"remote-add",
		&RemoteAdd{},
	)
}

type RemoteAdd struct {
	command_components.LocalWorkingCopy
	command_components.RemoteTransfer

	complete command_components.Complete

	proto sku.Proto
}

func (cmd *RemoteAdd) SetFlagSet(flagSet *flag.FlagSet) {
	cmd.RemoteTransfer.SetFlagSet(flagSet)

	flagSet.Var(
		cmd.complete.GetFlagValueMetadataTags(&cmd.proto.Metadata),
		"tags",
		"tags added for new objects in `checkin`, `new`, `organize`",
	)

	cmd.proto.SetFlagSetDescription(
		flagSet,
		"description to use for the new repo",
	)
}

func (cmd RemoteAdd) Run(req command.Request) {
	local := cmd.MakeLocalWorkingCopy(req)
	req.Must(local.Lock)
	cmd.AddRemote(req, local, cmd.proto)
	req.Must(local.Unlock)
}
