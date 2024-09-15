package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
)

func (u *Env) Lock() (err error) {
	ui.Log().Caller(1, "Umwelt Lock")
	if err = u.fs_home.GetLockSmith().Lock(); err != nil {
		ui.Log().Caller(1, "Umwelt Lock Failure")
		err = errors.Wrap(err)
		return
	}

	ui.Log().Caller(1, "Umwelt Lock Success")
	return
}

// TODO print organize files that were created if dry run or make it possible to
// commit dry-run transactions
func (u *Env) Unlock() (err error) {
	ptl := u.PrinterTransactedLike()

	if u.storesInitialized {
		ui.Log().Printf("konfig has changes: %t", u.GetConfig().HasChanges())
		ui.Log().Printf("schlummernd has changes: %t", u.GetDormantIndex().HasChanges())

		var changes []string
		changes = append(changes, u.GetConfig().GetChanges()...)
		changes = append(changes, u.GetDormantIndex().GetChanges()...)
		u.GetStore().GetStreamIndex().SetNeedsFlushHistory(changes)

		ui.Log().Print("will flush bestandsaufnahme")
		if err = u.store.FlushBestandsaufnahme(ptl); err != nil {
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
			u.GetFSHome(),
			u.GetStore().GetBlobStore().GetTypeV0(),
			u.PrinterHeader(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		ui.Log().Print("will flush schlummernd")
		if err = u.dormantIndex.Flush(
			u.GetFSHome(),
			u.PrinterHeader(),
			u.config.DryRun,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		u.GetStore().GetStreamIndex().SetNeedsFlushHistory(changes)

		wg := iter.MakeErrorWaitGroupParallel()
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
	if err = u.fs_home.GetLockSmith().Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
