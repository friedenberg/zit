package store_with_lock

import (
	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/delta/age"
	"github.com/friedenberg/zit/src/delta/file_lock"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/foxtrot/akten"
	"github.com/friedenberg/zit/src/golf/hinweisen"
	store_checkout "github.com/friedenberg/zit/src/hotel/store_checkout"
	store_objekten "github.com/friedenberg/zit/src/hotel/store_objekten"
)

type Store struct {
	*umwelt.Umwelt
	lock           *file_lock.Lock
	zettels        *store_objekten.Store
	akten          akten.Akten
	age            age.Age
	store_checkout *store_checkout.Store
}

func New(u *umwelt.Umwelt) (s Store, err error) {
	s.Umwelt = u
	s.lock = file_lock.New(u.DirZit("Lock"))

	if err = s.lock.Lock(); err != nil {
		err = errors.Error(err)
		return
	}

	if s.age, err = u.Age(); err != nil {
		err = errors.Error(err)
		return
	}

	s.zettels = &store_objekten.Store{}

	if err = s.zettels.Initialize(u); err != nil {
		err = errors.Wrapped(err, "failed to initialize zettel meta store")
		return
	}

	if s.akten, err = akten.New(u.DirZit()); err != nil {
		err = errors.Error(err)
		return
	}

	csk := store_checkout.Konfig{
		Konfig:       u.Konfig,
		CacheEnabled: u.Konfig.CheckoutCacheEnabled,
	}

	logz.Print("initing checkout store")
	if s.store_checkout, err = store_checkout.New(csk, u.Cwd(), s.zettels); err != nil {
		logz.Print(err)
		err = errors.Error(err)
		return
	}

	logz.Print("done initing checkout store")

	return
}

func (s Store) Age() age.Age {
	return s.age
}

func (s Store) Zettels() *store_objekten.Store {
	return s.zettels
}

func (s Store) Hinweisen() hinweisen.Hinweisen {
	return s.zettels.Hinweisen()
}

func (s Store) Akten() akten.Akten {
	return s.akten
}

func (s Store) CheckoutStore() *store_checkout.Store {
	return s.store_checkout
}

func (s Store) Flush() (err error) {
	if err = s.Zettels().Flush(); err != nil {
		stdprinter.Err(err)
		err = errors.Error(err)
		return
	}

	if err = s.CheckoutStore().Flush(); err != nil {
		stdprinter.Err(err)
		err = errors.Error(err)
		return
	}

	//explicitly do not unlock if there was an error to encourage user interaction
	//and manual recovery
	if err = s.lock.Unlock(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
