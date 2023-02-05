package umwelt

import "github.com/friedenberg/zit/src/alfa/errors"

func (u *Umwelt) Lock() (err error) {
	errors.Log().Caller(1, "Umwelt Lock")
	if err = u.lock.Lock(); err != nil {
		errors.Log().Caller(1, "Umwelt Lock Failure")
		err = errors.Wrap(err)
		return
	}

	errors.Log().Caller(1, "Umwelt Lock Success")
	return
}

func (u *Umwelt) Unlock() (err error) {
	if u.storesInitialized {
		if err = u.storeUtil.Flush(); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = u.storeObjekten.Flush(); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = u.storeWorkingDirectory.Flush(); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = u.KonfigPtr().Flush(u.Standort()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	//explicitly do not unlock if there was an error to encourage user interaction
	//and manual recovery
	if err = u.lock.Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
