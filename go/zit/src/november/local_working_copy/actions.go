package local_working_copy

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
)

func (u *Repo) DeleteFiles(fs interfaces.Collection[*fd.FD]) (err error) {
	deleteOp := store_fs.DeleteCheckout{}

	if err = deleteOp.Run(
		u.GetConfig().GetCLIConfig().IsDryRun(),
		u.GetEnvRepo(),
		u.PrinterFDDeleted(),
		fs,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
