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

	for listOrError := range remoteInventoryListStore.AllInventoryLists() {
		if listOrError.Error != nil {
			req.CancelWithError(listOrError.Error)
			return
		}

		if err := localInventoryListStore.ImportInventoryList(
			remote.GetBlobStore(),
			listOrError.Element,
		); err != nil {
			req.CancelWithError(err)
			return
		}

	}
}
