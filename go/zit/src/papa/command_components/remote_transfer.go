package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
)

type RemoteTransfer struct {
	Remote
	repo.RemoteTransferOptions
}

func (cmd *RemoteTransfer) SetFlagSet(f *flag.FlagSet) {
	cmd.Remote.SetFlagSet(f)
	cmd.RemoteTransferOptions.SetFlagSet(f)
}

func (cmd *RemoteTransfer) PushAllToArchive(
	req command.Request,
	local, remote repo.Repo,
) {
	remoteInventoryListStore := remote.GetInventoryListStore()
	localInventoryListStore := local.GetInventoryListStore()

	for list, err := range localInventoryListStore.IterAllInventoryLists() {
		if err != nil {
			req.CancelWithError(err)
			return
		}

		if err := remoteInventoryListStore.ImportInventoryList(
			local.GetBlobStore(),
			list,
		); err != nil {
			req.CancelWithError(err)
			return
		}
	}
}
