package umwelt

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/iter"
)

func (u *Umwelt) Lock() (err error) {
	errors.Log().Caller(1, "Umwelt Lock")
	if err = u.standort.GetLockSmith().Lock(); err != nil {
		errors.Log().Caller(1, "Umwelt Lock Failure")
		err = errors.Wrap(err)
		return
	}

	errors.Log().Caller(1, "Umwelt Lock Success")
	return
}

func (u *Umwelt) Unlock() (err error) {
	ptl := u.PrinterTransactedLike()

	if u.storesInitialized {
		if err = u.store.FlushBestandsaufnahme(ptl); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = u.store.Flush(u.PrinterHeader()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = u.Konfig().Flush(
			u.Standort(),
			u.GetStore().GetAkten().GetTypV0(),
			u.PrinterHeader(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		wg := iter.MakeErrorWaitGroupParallel()
		wg.Do(
			func() error {
				// second store flush is necessary because of konfig changes
				return u.store.Flush(u.PrinterHeader())
			},
		)

		for _, vs := range u.virtualStores {
			wg.Do(vs.Flush)
		}

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
