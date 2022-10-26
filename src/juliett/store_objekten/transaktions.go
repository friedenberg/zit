package store_objekten

import (
	"io"
	"path"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/id"
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/delta/objekte"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/echo/transaktion"
	"github.com/friedenberg/zit/src/echo/zettel_stored"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

func (s Store) ReadLastTransaktion() (t transaktion.Transaktion, err error) {
	var all []transaktion.Transaktion

	if all, err = s.ReadAllTransaktions(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if len(all) == 0 {
		err = ErrNotFound{}
	}

	sort.Slice(all, func(i, j int) bool { return all[j].Time.Less(all[i].Time) })

	t = all[0]

	return
}

func (s Store) ReadAllTransaktions() (out []transaktion.Transaktion, err error) {
	var headNames []string

	d := s.standort.DirObjektenTransaktion()

	if headNames, err = files.ReadDirNames(d); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, hn := range headNames {
		errors.Print(hn)

		var tailNames []string

		if tailNames, err = files.ReadDirNames(d, hn); err != nil {
			err = errors.Wrap(err)
			return
		}

		for _, tn := range tailNames {
			errors.Print(tn)

			p := path.Join(d, hn, tn)

			var t transaktion.Transaktion

			if t, err = s.readTransaktion(p); err != nil {
				err = errors.Wrap(err)
				return
			}

			out = append(out, t)
		}
	}

	errors.Print("sorting")
	sort.Slice(out, func(i, j int) bool { return out[i].Time.Less(out[j].Time) })
	errors.Print("done")

	return
}

func (s Store) ReadTransaktion(t ts.Time) (tr transaktion.Transaktion, err error) {
	return s.readTransaktion(id.Path(t, s.standort.DirObjektenTransaktion()))
}

func (s Store) readTransaktion(p string) (t transaktion.Transaktion, err error) {
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

	t = tr.Transaktion

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

	if _, err = f.ReadFrom(&sz.Zettel, or); err != nil {
		err = errors.Wrapf(err, "%s", sh)
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
) (tz zettel_transacted.Zettel, err error) {
	tz.Named = z
	tz.Kopf = t.Time
	tz.Schwanz = t.Time

	var previous zettel_transacted.Zettel

	if previous, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(z.Hinweis); err == nil {
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
	t transaktion.Transaktion,
	o objekte.ObjekteWithIndex,
) (tz zettel_transacted.Zettel, err error) {
	ok := false

	var h *hinweis.Hinweis

	if h, ok = o.Id.(*hinweis.Hinweis); !ok {
		err = errors.Wrapf(err, "transacktion.Objekte Id was not hinweis but was %s", o.Id)
		return
	}

	tz.Named.Hinweis = *h

	if tz.Named.Stored, err = s.storedZettelFromSha(o.Sha); err != nil {
		err = errors.Wrapf(err, "failed to read zettel objekte: %s", tz.Named.Hinweis)
		return
	}

	if tz, err = s.transactedWithHead(tz.Named, t); err != nil {
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
	errors.Printf("adding zettel to transaktion: %s", z.Hinweis)

	if tz, err = s.transactedWithHead(z, s.Transaktion); err != nil {
		err = errors.Wrap(err)
		return
	}

	var mutter [2]ts.Time

	mutter[0] = tz.Mutter

	i := s.Transaktion.Add(
		objekte.Objekte{
			Gattung: gattung.Zettel,
			Mutter:  mutter,
			Id:      &z.Hinweis,
			Sha:     z.Stored.Sha,
		},
	)

	tz.TransaktionIndex.SetInt(i)

	return
}
