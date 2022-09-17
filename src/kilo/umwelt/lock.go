package umwelt

import "github.com/friedenberg/zit/src/alfa/errors"

func (u *Umwelt) Lock() (err error) {
	errors.Caller(1, "Umwelt Lock")
	if err = u.lock.Lock(); err != nil {
		errors.Caller(1, "Umwelt Lock Failure")
		err = errors.Wrap(err)
		return
	}

	errors.Caller(1, "Umwelt Lock Success")
	return
}

func (u *Umwelt) Unlock() (err error) {
	errors.Caller(1, "Umwelt Unlock")
	if u.storesInitialized {
		if err = u.storeObjekten.Flush(); err != nil {
			errors.Caller(1, "Umwelt Unlock Failure")
			errors.PrintErr(err)
			err = errors.Wrap(err)
			return
		}

		if err = u.storeWorkingDirectory.Flush(); err != nil {
			errors.Caller(1, "Umwelt Unlock Failure")
			errors.PrintErr(err)
			err = errors.Wrap(err)
			return
		}
	}

	//explicitly do not unlock if there was an error to encourage user interaction
	//and manual recovery
	if err = u.lock.Unlock(); err != nil {
		errors.Caller(1, "Umwelt Unlock Failure")
		err = errors.Wrap(err)
		return
	}

	errors.Caller(1, "Umwelt Unlock Success")

	return
}
