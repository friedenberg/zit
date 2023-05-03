package store_util

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
)

type TransaktionStore interface {
	TransaktionPath(kennung.Time) string
	ReadTransaktion(kennung.Time) (*transaktion.Transaktion, error)
	ReadLastTransaktion() (*transaktion.Transaktion, error)
	ReadAllTransaktions(schnittstellen.FuncIter[*transaktion.Transaktion]) error
}

func (s common) ReadLastTransaktion() (t *transaktion.Transaktion, err error) {
	if err = s.ReadAllTransaktions(
		collections.MakeSyncSerializer(
			func(t1 *transaktion.Transaktion) (err error) {
				if t != nil && t1.Time.Less(t.Time) {
					return
				}

				t = t1

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if t == nil {
		err = errors.Wrap(objekte_store.ErrNotFound{})
	}

	return
}

func (s common) ReadAllTransaktions(
	f schnittstellen.FuncIter[*transaktion.Transaktion],
) (err error) {
	if err = files.ReadDirNamesLevel2(
		func(p string) (err error) {
			var t *transaktion.Transaktion

			if t, err = s.readTransaktion(p); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = f(t); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
		s.GetStandort().DirObjektenTransaktion(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s common) TransaktionPath(t kennung.Time) (p string) {
	p = id.Path(t, s.GetStandort().DirObjektenTransaktion())

	return
}

func (s common) ReadTransaktion(t kennung.Time) (tr *transaktion.Transaktion, err error) {
	return s.readTransaktion(s.TransaktionPath(t))
}

func (s common) readTransaktion(p string) (t *transaktion.Transaktion, err error) {
	tr := &transaktion.Reader{}

	var or io.ReadCloser

	if or, err = s.ReadCloserObjekten(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer or.Close()

	if _, err = tr.ReadFrom(or); err != nil {
		err = errors.Wrapf(err, "file: '%s'", p)
		return
	}

	t = &tr.Transaktion

	return
}
