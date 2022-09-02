package store_objekten

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/charlie/zk_types"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/echo/transaktion"
	"github.com/friedenberg/zit/src/golf/zettel_formats"
	"github.com/friedenberg/zit/src/golf/zettel_stored"
	"github.com/friedenberg/zit/zettel_named"
	"github.com/friedenberg/zit/zettel_transacted"
)

func (s Store) ReadTransaktion(t ts.Time) (tr transaktion.Transaktion, err error) {
	return s.readTransaktion(id.Path(t, s.Umwelt.DirObjektenTransaktion()))
}

func (s Store) readTransaktion(p string) (t transaktion.Transaktion, err error) {
	tr := &transaktion.Reader{}

	var or io.ReadCloser

	if or, err = s.ReadCloserObjekten(p); err != nil {
		err = errors.Error(err)
		return
	}

	defer or.Close()

	if _, err = tr.ReadFrom(or); err != nil {
		err = errors.Error(err)
		return
	}

	t = tr.Transaktion

	return
}

func (s Store) storedZettelFromSha(sh sha.Sha) (sz zettel_stored.Stored, err error) {
	var or io.ReadCloser

	if or, err = s.ReadCloserObjekten(id.Path(sh, s.Umwelt.DirObjektenZettelen())); err != nil {
		err = ErrNotFound{Id: sh}
		return
	}

	defer or.Close()

	f := zettel_formats.Objekte{}

	if _, err = f.ReadFrom(&sz.Zettel, or); err != nil {
		err = errors.Error(err)
		return
	}

	sz.Sha = sh

	return
}

// should only be called when moving forward through time, as there is a
// dependency on the index being accurate for the immediate mutter of the zettel
// in the arguments
func (s *Store) transactedWithHead(
	z zettel_named.Zettel,
	t transaktion.Transaktion,
) (tz zettel_transacted.Transacted, err error) {
	tz.Named = z
	tz.Kopf = t.Time
	tz.Schwanz = t.Time

	var previous zettel_transacted.Transacted

	if previous, err = s.indexZettelenTails.Read(z.Hinweis); err == nil {
		tz.Mutter = previous.Schwanz
		tz.Kopf = previous.Kopf
	} else {
		if errors.Is(err, ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Error(err)
			return
		}
	}

	return
}

func (s Store) transactedZettelFromTransaktionObjekte(
	t transaktion.Transaktion,
	o transaktion.Objekte,
) (tz zettel_transacted.Transacted, err error) {
	ok := false

	var h *hinweis.Hinweis

	if h, ok = o.Id.(*hinweis.Hinweis); !ok {
		err = errors.Wrapped(err, "transacktion.Objekte Id was not hinweis but was %s", o.Id)
		return
	}

	tz.Named.Hinweis = *h

	if tz.Named.Stored, err = s.storedZettelFromSha(o.Sha); err != nil {
		err = errors.Wrapped(err, "failed to find zettel objekte for hinweis: %s", tz.Named.Hinweis)
		return
	}

	if tz, err = s.transactedWithHead(tz.Named, t); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s Store) writeTransaktion() (err error) {
	if len(s.Transaktion.Objekten) == 0 {
		logz.Print("not writing Transaktion as there aren't any Objekten")
		return
	}

	logz.Printf("writing Transaktion with %d Objekten", len(s.Transaktion.Objekten))

	var p string

	if p, err = id.MakeDirIfNecessary(s.Transaktion.Time, s.Umwelt.DirObjektenTransaktion()); err != nil {
		err = errors.Error(err)
		return
	}

	var w io.WriteCloser

	if w, err = s.WriteCloserObjekten(p); err != nil {
		err = errors.Error(err)
		return
	}

	defer w.Close()

	f := transaktion.Writer{Transaktion: s.Transaktion}

	if _, err = f.WriteTo(w); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s *Store) addZettelToTransaktion(z zettel_named.Zettel) (tz zettel_transacted.Transacted, err error) {
	logz.Printf("adding zettel to transaktion: %s", z.Hinweis)

	if tz, err = s.transactedWithHead(z, s.Transaktion); err != nil {
		err = errors.Error(err)
		return
	}

	var mutter [2]ts.Time

	mutter[0] = tz.Mutter

	s.Transaktion.Objekten = append(
		s.Transaktion.Objekten,
		transaktion.Objekte{
			Type:   zk_types.TypeZettel,
			Mutter: mutter,
			Id:     &z.Hinweis,
			Sha:    z.Stored.Sha,
		},
	)

	return
}
