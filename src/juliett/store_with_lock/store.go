package store_with_lock

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/file_lock"
	"github.com/friedenberg/zit/src/delta/umwelt"
	"github.com/friedenberg/zit/src/echo/akten"
	"github.com/friedenberg/zit/src/hotel/store_objekten"
	"github.com/friedenberg/zit/src/india/store_working_directory"
)

type Store struct {
	*umwelt.Umwelt
	lock                  *file_lock.Lock
	storeObjekten         *store_objekten.Store
	akten                 akten.Akten
	age                   age.Age
	storeWorkingDirectory *store_working_directory.Store
}

func New(u *umwelt.Umwelt) (s Store, err error) {
	s.Umwelt = u
	s.lock = file_lock.New(u.DirZit("Lock"))

	if err = s.lock.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.age, err = u.Age(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.storeObjekten, err = store_objekten.Make(s.age, u.Konfig, u.Standort); err != nil {
		err = errors.Wrapf(err, "failed to initialize zettel meta store")
		return
	}

	if s.akten, err = akten.New(u.DirZit()); err != nil {
		err = errors.Wrap(err)
		return
	}

	csk := store_working_directory.Konfig{
		Konfig:       u.Konfig,
		CacheEnabled: u.Konfig.CheckoutCacheEnabled,
	}

	errors.Print("initing checkout store")
	if s.storeWorkingDirectory, err = store_working_directory.New(csk, u.Cwd(), s.storeObjekten); err != nil {
		errors.Print(err)
		err = errors.Wrap(err)
		return
	}

	errors.Print("done initing checkout store")

	return
}

func (s Store) Age() age.Age {
	return s.age
}

func (s Store) StoreObjekten() *store_objekten.Store {
	return s.storeObjekten
}

func (s Store) Akten() akten.Akten {
	return s.akten
}

func (s Store) StoreWorkingDirectory() *store_working_directory.Store {
	return s.storeWorkingDirectory
}

func (s Store) Flush() (err error) {
	if err = s.StoreObjekten().Flush(); err != nil {
		errors.PrintErr(err)
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreWorkingDirectory().Flush(); err != nil {
		errors.PrintErr(err)
		err = errors.Wrap(err)
		return
	}

	//explicitly do not unlock if there was an error to encourage user interaction
	//and manual recovery
	if err = s.lock.Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
