package umwelt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
)

func (u *Umwelt) DeleteFiles(fs schnittstellen.Iterable[*fd.FD]) (err error) {
	deleteOp := store_fs.DeleteCheckout{}

	if err = deleteOp.Run(
		u.GetKonfig().DryRun,
		u.Standort(),
		u.PrinterFDDeleted(),
		fs,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
