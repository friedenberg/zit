package store_objekten

import (
	"io"
	"reflect"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/hinweisen"
	"github.com/friedenberg/zit/src/foxtrot/id"
	"github.com/friedenberg/zit/src/golf/age_io"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/lima/store_verzeichnisse"
)

type zettelStore struct {
	common *common

	protoZettel zettel.ProtoZettel

	zettelTransactedWriter zettel.TransactedWriters

	*indexKennung
	hinweisen *hinweisen.Hinweisen
	*indexEtiketten

	verzeichnisseSchwanzen *verzeichnisseSchwanzen
	verzeichnisseAll       *store_verzeichnisse.Zettelen

	objekte.TransactedInflator[
		zettel.Objekte,
		*zettel.Objekte,
		hinweis.Hinweis,
		*hinweis.Hinweis,
		zettel.Verzeichnisse2,
		*zettel.Verzeichnisse2,
	]

	pool *zettel.PoolVerzeichnisse
}

func makeZettelStore(
	sa *common,
	p *zettel.PoolVerzeichnisse,
) (s *zettelStore, err error) {
	s = &zettelStore{
		common:      sa,
		pool:        p,
		protoZettel: zettel.MakeProtoZettel(sa.Konfig()),
		TransactedInflator: objekte.MakeTransactedInflator[
			zettel.Objekte,
			*zettel.Objekte,
			hinweis.Hinweis,
			*hinweis.Hinweis,
			zettel.Verzeichnisse2,
			*zettel.Verzeichnisse2,
		](
			sa,
			func(sh sha.Sha) (r sha.ReadCloser, err error) {
				return sa.ReadCloserObjekten(
					id.Path(sh, sa.Standort.DirObjektenZettelen()),
				)
			},
			&zettel.FormatObjekte{
				IgnoreTypErrors: true,
			},
			nil,
		),
	}

	if s.hinweisen, err = hinweisen.New(s.common.Standort.DirZit()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.verzeichnisseSchwanzen, err = makeVerzeichnisseSchwanzen(
		s.common,
		p,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.verzeichnisseAll, err = store_verzeichnisse.MakeZettelen(
		s.common.Konfig(),
		s.common.Standort.DirVerzeichnisseZettelenNeue(),
		s.common,
		p,
		nil,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.indexKennung, err = newIndexKennung(
		s.common.Konfig(),
		s.common,
		s.hinweisen,
		s.common.Standort.DirVerzeichnisse("Kennung"),
	); err != nil {
		err = errors.Wrapf(err, "failed to init kennung index")
		return
	}

	if s.indexEtiketten, err = newIndexEtiketten(
		s.common.Standort.FileVerzeichnisseEtiketten(),
		s.common,
	); err != nil {
		err = errors.Wrapf(err, "failed to init zettel index")
		return
	}

	return
}

func (s *zettelStore) Hinweisen() *hinweisen.Hinweisen {
	return s.hinweisen
}

func (s *zettelStore) Flush() (err error) {
	if err = s.verzeichnisseSchwanzen.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.verzeichnisseAll.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.indexEtiketten.Flush(); err != nil {
		err = errors.Wrapf(err, "failed to flush new zettel index")
		return
	}

	if err = s.indexKennung.Flush(); err != nil {
		err = errors.Wrapf(err, "failed to flush new kennung index")
		return
	}

	return
}

func (s *zettelStore) SetZettelTransactedLogWriter(
	ztlw zettel.TransactedWriters,
) {
	s.zettelTransactedWriter = ztlw
}

func (s zettelStore) WriteZettelObjekte(z zettel.Objekte) (sh sha.Sha, err error) {
	//no lock required

	var w *age_io.Mover

	mo := age_io.MoveOptions{
		Age:                      s.common.Age,
		FinalPath:                s.common.Standort.DirObjektenZettelen(),
		GenerateFinalPathFromSha: true,
	}

	if w, err = age_io.NewMover(mo); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, w.Close)

	c := zettel.ObjekteFormatterContext{
		Zettel: z,
	}

	f := zettel.FormatObjekte{}

	if _, err = f.Format(w, &c.Zettel); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = w.Sha()

	return
}

func (s *zettelStore) writeNamedZettelToIndex(
	tz *zettel.Transacted,
) (err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "write named zettel to index",
		}

		return
	}

	errors.Log().Printf("writing zettel to index: %s", tz.Sku)

	if err = s.verzeichnisseSchwanzen.Add(tz, tz.Sku.Kennung.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.verzeichnisseAll.Add(tz, tz.Sku.Sha.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.indexKennung.addHinweis(tz.Sku.Kennung); err != nil {
		if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
			errors.Log().Printf("kennung does not contain value: %s", err)
			err = nil
		} else {
			err = errors.Wrapf(err, "failed to write zettel to index: %s", tz.Sku)
			return
		}
	}

	if err = s.common.Abbr.addZettelTransacted(tz); err != nil {
		err = errors.Wrapf(err, "failed to write zettel to index: %s", tz.Sku)
		return
	}

	return
}

func (s zettelStore) ReadHinweisSchwanzen(
	h hinweis.Hinweis,
) (zv *zettel.Transacted, err error) {
	return s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(h)
}

func (i *zettelStore) ReadAllSchwanzenVerzeichnisse(
	w collections.WriterFunc[*zettel.Verzeichnisse],
) (err error) {
	return i.verzeichnisseSchwanzen.ReadMany(w)
}

func (i *zettelStore) ReadAllVerzeichnisse(
	w collections.WriterFunc[*zettel.Verzeichnisse],
) (err error) {
	return i.verzeichnisseAll.ReadMany(w)
}

func (s zettelStore) ReadOne(
	i id.Id,
) (tz *zettel.Transacted, err error) {
	switch tid := i.(type) {
	case hinweis.Hinweis:
		if tz, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(tid); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = errors.Errorf("unsupported identifier: %s, %#v", i, reflect.ValueOf(i))
	}

	return
}

func (s *zettelStore) Create(
	in zettel.Objekte,
) (tz *zettel.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "create",
		}

		return
	}

	if in.IsEmpty() || s.protoZettel.Equals(in) {
		err = errors.Normalf("zettel is empty")
		return
	}

	s.protoZettel.Apply(&in)

	if err = in.ApplyKonfig(s.common.Konfig()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var shaObj sha.Sha

	if shaObj, err = s.WriteZettelObjekte(in); err != nil {
		err = errors.Wrap(err)
		return
	}

	//If the zettel exists, short circuit and return that
	if tz2, err2 := s.ReadOne(shaObj); err2 == nil {
		tz = tz2
		return
	}

	var ken hinweis.Hinweis

	if ken, err = s.indexKennung.createHinweis(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz, err = s.addZettelToTransaktion(
		&in,
		&shaObj,
		&ken,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.indexEtiketten.add(tz.Objekte.Etiketten); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO-P2 assert no changes
	if err = s.zettelTransactedWriter.New(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO support dry run
func (s *zettelStore) Update(
	z *zettel.Objekte,
	h *hinweis.Hinweis,
) (tz *zettel.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "update",
		}

		return
	}

	if err = z.ApplyKonfig(s.common.Konfig()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var mutter *zettel.Transacted

	if mutter, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(
		*h,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var shaObj sha.Sha

	if shaObj, err = s.WriteZettelObjekte(*z); err != nil {
		err = errors.Wrap(err)
		return
	}

	if shaObj.Equals(mutter.Sku.Sha) {
		tz = mutter

		if err = s.zettelTransactedWriter.Unchanged(tz); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if tz, err = s.addZettelToTransaktion(
		z,
		&shaObj,
		h,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.indexEtiketten.addZettelWithOptionalMutter(tz, mutter); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.zettelTransactedWriter.Updated(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s zettelStore) AllInChain(h hinweis.Hinweis) (c []*zettel.Verzeichnisse, err error) {
	mst := zettel.MakeMutableSetUnique(0)

	if err = s.verzeichnisseAll.ReadMany(
		func(z *zettel.Verzeichnisse) (err error) {
			if !z.Transacted.Sku.Kennung.Equals(&h) {
				err = io.EOF
				return
			}

			return
		},
		mst.AddAndDoNotRepool,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	c = mst.Elements()

	sort.Slice(
		c,
		func(i, j int) bool { return c[i].Transacted.Sku.Less(&c[j].Transacted.Sku) },
	)

	return
}

func (s *zettelStore) addZettelToTransaktion(
	zo *zettel.Objekte,
	zs *sha.Sha,
	zk *hinweis.Hinweis,
) (tz *zettel.Transacted, err error) {
	errors.Log().Printf("adding zettel to transaktion: %s", zk)

	if tz, err = s.transactedWithHead(
		*zo,
		*zk,
		&s.common.Transaktion,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz.Sku.Kennung = *zk
	tz.Sku.Sha = *zs

	s.common.Transaktion.Add2(&tz.Sku)

	return
}

// TODO-P1 modify this to not require immediate mutter
// should only be called when moving forward through time, as there is a
// dependency on the index being accurate for the immediate mutter of the zettel
// in the arguments
func (s *zettelStore) transactedWithHead(
	z zettel.Objekte,
	h hinweis.Hinweis,
	t *transaktion.Transaktion,
) (tz *zettel.Transacted, err error) {
	tz = &zettel.Transacted{
		Objekte: z,
		Sku: sku.Transacted[hinweis.Hinweis, *hinweis.Hinweis]{
			Kennung: h,
			Kopf:    t.Time,
			Schwanz: t.Time,
		},
	}

	var previous *zettel.Transacted

	if previous, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(h); err == nil {
		tz.Sku.Mutter[0] = previous.Sku.Schwanz
		tz.Sku.Kopf = previous.Sku.Kopf
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

func (s *zettelStore) reindexOne(
	t *transaktion.Transaktion,
	o *sku.Sku,
) (err error) {
	var tz *zettel.Transacted

	if tz, err = s.Inflate(t, o); err != nil {
		//TODO-P2 decide on how to handle format errors
		errors.Err().Print(err)
		err = nil
		// err = errors.Wrap(err)
		return
	}

	var mutter *zettel.Transacted

	if mutter1, err := s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(tz.Sku.Kennung); err == nil {
		mutter = mutter1
	}

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if mutter == nil {
		if err = s.zettelTransactedWriter.New(tz); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = s.zettelTransactedWriter.Updated(tz); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.indexEtiketten.addZettelWithOptionalMutter(tz, mutter); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
