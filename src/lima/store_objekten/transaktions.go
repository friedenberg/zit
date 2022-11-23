package store_objekten

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/golf/zettel_stored"
	"github.com/friedenberg/zit/src/hotel/zettel_named"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
)

func (s Store) ReadLastTransaktion() (t *transaktion.Transaktion, err error) {
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
		err = ErrNotFound{}
	}

	return
}

func (s Store) ReadAllTransaktions(
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
		s.standort.DirObjektenTransaktion(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Store) ReadTransaktion(t ts.Time) (tr *transaktion.Transaktion, err error) {
	return s.readTransaktion(id.Path(t, s.standort.DirObjektenTransaktion()))
}

func (s Store) readTransaktion(p string) (t *transaktion.Transaktion, err error) {
	tr := &transaktion.Reader{}

	var or io.ReadCloser

	if or, err = s.ReadCloserObjekten(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer or.Close()

	if _, err = tr.ReadFrom(or); err != nil {
		err = errors.Wrap(err)
		return
	}

	t = &tr.Transaktion

	return
}

func (s Store) storedZettelFromSha(sh sha.Sha) (sz zettel_stored.Stored, err error) {
	var or io.ReadCloser

	if or, err = s.ReadCloserObjekten(id.Path(sh, s.standort.DirObjektenZettelen())); err != nil {
		err = ErrNotFound{Id: sh}
		return
	}

	defer or.Close()

	f := zettel.Objekte{
		IgnoreTypErrors: true,
	}

	c := zettel.FormatContextRead{
		In: or,
	}

	if _, err = f.ReadFrom(&c); err != nil {
		err = errors.Wrapf(err, "%s", sh)
		return
	}

	sz.Objekte = c.Zettel
	sz.Sha = sh

	return
}

// should only be called when moving forward through time, as there is a
// dependency on the index being accurate for the immediate mutter of the zettel
// in the arguments
func (s *Store) transactedWithHead(
	z zettel_named.Zettel,
	t transaktion.Transaktion,
) (tz zettel_transacted.Zettel, err error) {
	tz.Named = z
	tz.Kopf = t.Time
	tz.Schwanz = t.Time

	var previous zettel_transacted.Zettel

	if previous, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(z.Kennung); err == nil {
		tz.Mutter = previous.Schwanz
		tz.Kopf = previous.Kopf
	} else {
		if errors.Is(err, ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s Store) transactedZettelFromTransaktionObjekte(
	t *transaktion.Transaktion,
	o *sku.Indexed,
) (tz zettel_transacted.Zettel, err error) {
	ok := false

	var h *hinweis.Hinweis

	if h, ok = o.Id.(*hinweis.Hinweis); !ok {
		err = errors.Wrapf(err, "transacktion.Objekte Id was not hinweis but was %s", o.Id)
		return
	}

	tz.Named.Kennung = *h

	if tz.Named.Stored, err = s.storedZettelFromSha(o.Sha); err != nil {
		err = errors.Wrapf(err, "failed to read zettel objekte: %s", tz.Named.Kennung)
		return
	}

	if tz, err = s.transactedWithHead(tz.Named, *t); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz.TransaktionIndex = o.Index

	return
}

func (s Store) writeTransaktion() (err error) {
	if s.Transaktion.Len() == 0 {
		errors.Print("not writing Transaktion as there aren't any Objekten")
		return
	}

	errors.Printf("writing Transaktion with %d Objekten", s.Transaktion.Len())

	var p string

	if p, err = id.MakeDirIfNecessary(s.Transaktion.Time, s.standort.DirObjektenTransaktion()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var w io.WriteCloser

	if w, err = s.WriteCloserObjekten(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer w.Close()

	f := transaktion.Writer{Transaktion: s.Transaktion}

	if _, err = f.WriteTo(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) addZettelToTransaktion(z zettel_named.Zettel) (tz zettel_transacted.Zettel, err error) {
	errors.Printf("adding zettel to transaktion: %s", z.Kennung)

	if tz, err = s.transactedWithHead(z, s.Transaktion); err != nil {
		err = errors.Wrap(err)
		return
	}

	var mutter [2]ts.Time

	mutter[0] = tz.Mutter

	i := s.Transaktion.Add(
		sku.Sku{
			Gattung: gattung.Zettel,
			Mutter:  mutter,
			Id:      &z.Kennung,
			Sha:     z.Stored.Sha,
		},
	)

	tz.TransaktionIndex.SetInt(i)

	return
}
