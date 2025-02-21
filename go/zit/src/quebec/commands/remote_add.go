package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
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
	_, remoteObject := cmd.CreateRemoteObject(req, local)

	var id ids.RepoId

	if err := id.Set(req.PopArg("repo-id")); err != nil {
		req.CancelWithError(err)
	}

	req.AssertNoMoreArgs()

	if err := remoteObject.ObjectId.SetWithIdLike(&id); err != nil {
		req.CancelWithError(err)
	}

	// TODO connect to remote and get public key and validate

	cmd.proto.Apply(remoteObject.GetMetadata(), genres.Repo)

	req.Must(local.Lock)

	if err := local.GetStore().CreateOrUpdate(
		remoteObject,
		sku.StoreOptions{
			ApplyProto: true,
		},
	); err != nil {
		req.CancelWithError(err)
	}

	req.Must(local.Unlock)
}
