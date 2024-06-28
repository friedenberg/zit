package umwelt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
)

func (u *Umwelt) Lock() (err error) {
	ui.Log().Caller(1, "Umwelt Lock")
	if err = u.standort.GetLockSmith().Lock(); err != nil {
		ui.Log().Caller(1, "Umwelt Lock Failure")
		err = errors.Wrap(err)
		return
	}

	ui.Log().Caller(1, "Umwelt Lock Success")
	return
}

func (u *Umwelt) Unlock() (err error) {
	ptl := u.PrinterTransactedLike()

	if u.storesInitialized {
		ui.Log().Printf("konfig has changes: %t", u.GetKonfig().HasChanges())
		ui.Log().Printf("schlummernd has changes: %t", u.Schlummernd().HasChanges())

		var changes []string
		changes = append(changes, u.GetKonfig().GetChanges()...)
		changes = append(changes, u.Schlummernd().GetChanges()...)
		u.GetStore().GetVerzeichnisse().SetNeedsFlushHistory(changes)

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
		if err = u.konfig.Flush(
			u.Standort(),
			u.GetStore().GetAkten().GetTypV0(),
			u.PrinterHeader(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		ui.Log().Print("will flush schlummernd")
		if err = u.schlummernd.Flush(
			u.Standort(),
			u.PrinterHeader(),
			u.konfig.DryRun,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		u.GetStore().GetVerzeichnisse().SetNeedsFlushHistory(changes)

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
	if err = u.standort.GetLockSmith().Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
