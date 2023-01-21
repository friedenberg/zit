package store_objekten

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/golf/transaktion"
)

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
		err = errors.Wrap(ErrNotFound{})
	}

	return
}

func (s common) ReadAllTransaktions(
	f collections.WriterFunc[*transaktion.Transaktion],
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

func (s common) TransaktionPath(t ts.Time) (p string) {
	p = id.Path(t, s.GetStandort().DirObjektenTransaktion())

	return
}

func (s common) ReadTransaktion(t ts.Time) (tr *transaktion.Transaktion, err error) {
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

func (s common) writeTransaktion() (err error) {
	if s.GetTransaktion().Skus.Len() == 0 {
		errors.Log().Print("not writing Transaktion as there aren't any Objekten")
		return
	}

	errors.Log().Printf(
		"writing Transaktion with %d Objekten",
		s.GetTransaktion().Skus.Len(),
	)

	var p string

	if p, err = id.MakeDirIfNecessary(
		s.GetTransaktion().Time,
		s.GetStandort().DirObjektenTransaktion(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var w io.WriteCloser

	if w, err = s.WriteCloserObjekten(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, w.Close)

	f := transaktion.Writer{Transaktion: *s.GetTransaktion()}

	if _, err = f.WriteTo(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
