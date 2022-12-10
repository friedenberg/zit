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

//TODO-P0 flush konfig_compiled
func (u *Umwelt) Unlock() (err error) {
	errors.Log().Caller(1, "Umwelt Unlock")
	if u.storesInitialized {
		if err = u.storeObjekten.Flush(); err != nil {
			errors.Log().Caller(1, "Umwelt Unlock Failure")
			errors.PrintErr(err)
			err = errors.Wrap(err)
			return
		}

		if err = u.storeWorkingDirectory.Flush(); err != nil {
			errors.Log().Caller(1, "Umwelt Unlock Failure")
			errors.PrintErr(err)
			err = errors.Wrap(err)
			return
		}
	}

	//explicitly do not unlock if there was an error to encourage user interaction
	//and manual recovery
	if err = u.lock.Unlock(); err != nil {
		errors.Log().Caller(1, "Umwelt Unlock Failure")
		err = errors.Wrap(err)
		return
	}

	errors.Log().Caller(1, "Umwelt Unlock Success")

	return
}
