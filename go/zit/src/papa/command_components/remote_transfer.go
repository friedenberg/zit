package command_components

import (
	"flag"

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
