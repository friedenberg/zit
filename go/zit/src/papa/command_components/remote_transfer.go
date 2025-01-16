package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
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
	local, remote repo.Archive,
) {
	remoteInventoryListStore := remote.GetInventoryListStore()
	localInventoryListStore := local.GetInventoryListStore()

	if err := remoteInventoryListStore.ReadAllInventoryLists(
		func(sk *sku.Transacted) (err error) {
			if err = localInventoryListStore.ImportInventoryList(
				remote.GetBlobStore(),
				sk,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		req.CancelWithError(err)
	}
}
