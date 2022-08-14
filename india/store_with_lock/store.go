package store_with_lock

import (
	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/stdprinter"
	"github.com/friedenberg/zit/delta/age"
	"github.com/friedenberg/zit/delta/file_lock"
	"github.com/friedenberg/zit/echo/umwelt"
	"github.com/friedenberg/zit/foxtrot/akten"
	"github.com/friedenberg/zit/golf/hinweisen"
	checkout_store "github.com/friedenberg/zit/golf/store_checkout"
	objekten "github.com/friedenberg/zit/golf/store_objekten"
)

type Store struct {
	*umwelt.Umwelt
	lock           *file_lock.Lock
	zettels        *objekten.Store
	akten          akten.Akten
	age            age.Age
	checkout_store *checkout_store.Store
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

	s.zettels = &objekten.Store{}

	if err = s.zettels.Initialize(u); err != nil {
		err = errors.Wrapped(err, "failed to initialize zettel meta store")
		return
	}

	if s.akten, err = akten.New(u.DirZit()); err != nil {
		err = errors.Error(err)
		return
	}

	csk := checkout_store.Konfig{
		CacheEnabled: u.Konfig.CheckoutCacheEnabled,
	}

	logz.Print("initing checkout store")
	if s.checkout_store, err = checkout_store.New(csk, u.Cwd(), s.zettels); err != nil {
		logz.Print(err)
		err = errors.Error(err)
		return
	}

	return
}

func (s Store) Age() age.Age {
	return s.age
}

func (s Store) Zettels() *objekten.Store {
	return s.zettels
}

func (s Store) Hinweisen() hinweisen.Hinweisen {
	return s.zettels.Hinweisen()
}

func (s Store) Akten() akten.Akten {
	return s.akten
}

func (s Store) CheckoutStore() *checkout_store.Store {
	return s.checkout_store
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
