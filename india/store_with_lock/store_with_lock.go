package store_with_lock

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/charlie/age"
	"github.com/friedenberg/zit/charlie/file_lock"
	"github.com/friedenberg/zit/golf/checkout_store"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/foxtrot/akten"
	"github.com/friedenberg/zit/foxtrot/etiketten"
	"github.com/friedenberg/zit/foxtrot/hinweisen"
	"github.com/friedenberg/zit/hotel/zettels"
)

type Store struct {
	*umwelt.Umwelt
	lock           *file_lock.Lock
	zettels        zettels.Zettels
	akten          akten.Akten
	age            age.Age
	checkout_store checkout_store.Store
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

	if s.zettels, err = zettels.New(u, s.age); err != nil {
		err = errors.Error(err)
		return
	}

	if s.akten, err = akten.New(u.DirZit()); err != nil {
		err = errors.Error(err)
		return
	}

	if s.checkout_store, err = checkout_store.New(u.Cwd(), s.zettels); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s Store) Age() age.Age {
	return s.age
}

func (s Store) Zettels() zettels.Zettels {
	return s.zettels
}

func (s Store) Hinweisen() hinweisen.Hinweisen {
	return s.zettels.Hinweisen()
}

func (s Store) Etiketten() etiketten.Etiketten {
	return s.zettels.Etiketten()
}

func (s Store) Akten() akten.Akten {
	return s.akten
}

func (s Store) CheckoutStore() checkout_store.Store {
	return s.checkout_store
}

func (s Store) Flush() (err error) {
	if err = s.zettels.Flush(); err != nil {
		err = errors.Error(err)
		return
	}

	//TODO always do this?
	if err = s.lock.Unlock(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
