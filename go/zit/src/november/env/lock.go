package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
)

func (u *Local) Lock() (err error) {
	ui.Log().Caller(1, "Umwelt Lock")
	if err = u.dirLayout.GetLockSmith().Lock(); err != nil {
		ui.Log().Caller(1, "Umwelt Lock Failure")
		err = errors.Wrap(err)
		return
	}

	ui.Log().Caller(1, "Umwelt Lock Success")
	return
}

// TODO print organize files that were created if dry run or make it possible to
// commit dry-run transactions
func (u *Local) Unlock() (err error) {
	ptl := u.PrinterTransacted()

	if u.storesInitialized {
		ui.Log().Printf("konfig has changes: %t", u.GetConfig().HasChanges())
		ui.Log().Printf("schlummernd has changes: %t", u.GetDormantIndex().HasChanges())

		var changes []string
		changes = append(changes, u.GetConfig().GetChanges()...)
		changes = append(changes, u.GetDormantIndex().GetChanges()...)
		u.GetStore().GetStreamIndex().SetNeedsFlushHistory(changes)

		ui.Log().Print("will flush bestandsaufnahme")
		if err = u.store.FlushInventoryList(ptl); err != nil {
			err = errors.Wrap(err)
			return
		}

		ui.Log().Print("will flush store")
		if err = u.store.Flush(
			u.PrinterHeader(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		ui.Log().Print("will flush konfig")
		if err = u.config.Flush(
			u.GetDirectoryLayout(),
			u.GetStore().GetBlobStore(),
			u.PrinterHeader(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		ui.Log().Print("will flush schlummernd")
		if err = u.dormantIndex.Flush(
			u.GetDirectoryLayout(),
			u.PrinterHeader(),
			u.config.DryRun,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		u.GetStore().GetStreamIndex().SetNeedsFlushHistory(changes)

		wg := quiter.MakeErrorWaitGroupParallel()
		wg.Do(
			func() error {
				ui.Log().Print("will flush store second time")
				// second store flush is necessary because of konfig changes
				return u.store.Flush(
					u.PrinterHeader(),
				)
			},
		)

		if err = wg.GetError(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	// explicitly do not unlock if there was an error to encourage user
	// interaction
	// and manual recovery
	if err = u.dirLayout.GetLockSmith().Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
