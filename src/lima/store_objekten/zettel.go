package store_objekten

import (
	"io"
	"reflect"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/hinweisen"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/echo/age_io"
	"github.com/friedenberg/zit/src/echo/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/india/zettel"
	"github.com/friedenberg/zit/src/juliett/zettel_verzeichnisse"
	"github.com/friedenberg/zit/src/kilo/store_verzeichnisse"
)

type zettelStore struct {
	common *common

	indexAbbr *indexAbbr

	protoZettel zettel.ProtoZettel

	zettelTransactedWriter ZettelTransactedLogWriters

	*indexKennung
	hinweisen *hinweisen.Hinweisen
	*indexEtiketten

	verzeichnisseSchwanzen *verzeichnisseSchwanzen
	verzeichnisseAll       *store_verzeichnisse.Zettelen

	pool zettel_verzeichnisse.Pool
}

func makeZettelStore(
	sa *common,
	p zettel_verzeichnisse.Pool,
	ia *indexAbbr,
) (s *zettelStore, err error) {
	s = &zettelStore{
		common:      sa,
		indexAbbr:   ia,
		pool:        p,
		protoZettel: zettel.MakeProtoZettel(),
	}

	if s.hinweisen, err = hinweisen.New(s.common.Standort.DirZit()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.protoZettel.Typ.Set(s.common.Konfig.Transacted.Objekte.Akte.DefaultTyp.Sku.Kennung.String()); err != nil {
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
		s.common.Konfig,
		s.common.Standort.DirVerzeichnisseZettelenNeue(),
		s.common,
		p,
		nil,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.indexKennung, err = newIndexKennung(
		s.common.Konfig,
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

// TODO add archived state
type ZettelTransactedLogWriters struct {
	New, Updated, Archived, Unchanged collections.WriterFunc[*zettel.Transacted]
}

func (s *zettelStore) SetZettelTransactedLogWriter(
	ztlw ZettelTransactedLogWriters,
) {
	s.zettelTransactedWriter = ztlw
}

func (s zettelStore) WriteZettelObjekte(z zettel.Zettel) (sh sha.Sha, err error) {
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

	defer w.Close()

	c := zettel.FormatContextWrite{
		Zettel: z,
		Out:    w,
	}

	f := zettel.Objekte{}

	if _, err = f.WriteTo(c); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = w.Sha()

	return
}

func (s *zettelStore) writeNamedZettelToIndex(tz zettel.Transacted) (err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "write named zettel to index",
		}

		return
	}

	errors.Log().Printf("writing zettel to index: %s", tz.Named)

	if err = s.verzeichnisseSchwanzen.Add(tz, tz.Named.Kennung.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.verzeichnisseAll.Add(tz, tz.Named.Stored.Sha.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.indexKennung.addHinweis(tz.Named.Kennung); err != nil {
		if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
			errors.Log().Printf("kennung does not contain value: %s", err)
			err = nil
		} else {
			err = errors.Wrapf(err, "failed to write zettel to index: %s", tz.Named)
			return
		}
	}

	if err = s.indexAbbr.addZettelTransacted(tz); err != nil {
		err = errors.Wrapf(err, "failed to write zettel to index: %s", tz.Named)
		return
	}

	return
}

func (s zettelStore) ReadHinweisSchwanzen(
	h hinweis.Hinweis,
) (zv zettel.Transacted, err error) {
	return s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(h)
}

func (i *zettelStore) ReadAllSchwanzenVerzeichnisse(
	ws ...collections.WriterFunc[*zettel_verzeichnisse.Zettel],
) (err error) {
	return i.verzeichnisseSchwanzen.ReadMany(ws...)
}

func (s zettelStore) ReadAllSchwanzenTransacted(
	ws ...collections.WriterFunc[*zettel.Transacted],
) (err error) {
	w := zettel_verzeichnisse.MakeWriterZettelTransacted(
		collections.MakeChain(ws...),
	)

	return s.ReadAllSchwanzenVerzeichnisse(w)
}

func (i *zettelStore) ReadAllVerzeichnisse(
	ws ...collections.WriterFunc[*zettel_verzeichnisse.Zettel],
) (err error) {
	return i.verzeichnisseAll.ReadMany(ws...)
}

func (s zettelStore) ReadAllTransacted(
	ws ...collections.WriterFunc[*zettel.Transacted],
) (err error) {
	w := zettel_verzeichnisse.MakeWriterZettelTransacted(
		collections.MakeChain(ws...),
	)

	return s.ReadAllVerzeichnisse(w)
}

func (s zettelStore) ReadOne(i id.Id) (tz zettel.Transacted, err error) {
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

func (s *zettelStore) Create(in zettel.Zettel) (tz zettel.Transacted, err error) {
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

	if err = in.ApplyKonfig(s.common.Konfig); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz.Named.Stored.Objekte = in

	if tz.Named.Stored.Sha, err = s.WriteZettelObjekte(tz.Named.Stored.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	//If the zettel exists, short circuit and return that
	if tz2, err2 := s.ReadOne(tz.Named.Stored.Sha); err2 == nil {
		tz = tz2
		return
	}

	if tz.Named.Kennung, err = s.indexKennung.createHinweis(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz, err = s.addZettelToTransaktion(tz.Named); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.indexEtiketten.add(tz.Named.Stored.Objekte.Etiketten); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO-P2 assert no changes
	if err = s.zettelTransactedWriter.New(&tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) CreateWithHinweis(
	in zettel.Zettel,
	h hinweis.Hinweis,
) (tz zettel.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "create with hinweis",
		}

		return
	}

	if in.IsEmpty() {
		err = errors.Normalf("zettel is empty")
		return
	}

	if err = in.ApplyKonfig(s.common.Konfig); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz.Named.Stored.Objekte = in

	if tz.Named.Stored.Sha, err = s.WriteZettelObjekte(tz.Named.Stored.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz, err = s.addZettelToTransaktion(tz.Named); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.indexEtiketten.add(tz.Named.Stored.Objekte.Etiketten); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.zettelTransactedWriter.New(&tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO support dry run
func (s *zettelStore) Update(
	z *zettel.Named,
) (tz zettel.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "update",
		}

		return
	}

	if err = z.Stored.Objekte.ApplyKonfig(s.common.Konfig); err != nil {
		err = errors.Wrap(err)
		return
	}

	var mutter zettel.Transacted

	if mutter, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(
		z.Kennung,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz.Named = *z

	if tz.Named.Stored.Sha, err = s.WriteZettelObjekte(z.Stored.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz.Named.Stored.Sha.Equals(mutter.Named.Stored.Sha) {
		tz = mutter

		if err = s.zettelTransactedWriter.Unchanged(&tz); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if tz, err = s.addZettelToTransaktion(tz.Named); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.indexEtiketten.addZettelWithOptionalMutter(&tz, &mutter); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.zettelTransactedWriter.Updated(&tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s zettelStore) AllInChain(h hinweis.Hinweis) (c []*zettel.Transacted, err error) {
	mst := zettel.MakeMutableSetUnique(0)

	if err = s.verzeichnisseAll.ReadMany(
		func(z *zettel_verzeichnisse.Zettel) (err error) {
			if !z.Transacted.Named.Kennung.Equals(&h) {
				err = io.EOF
				return
			}

			return
		},
		zettel_verzeichnisse.MakeWriterZettelTransacted(mst.AddAndDoNotRepool),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	c = mst.Elements()

	sort.Slice(
		c,
		func(i, j int) bool { return c[i].SkuTransacted().Less(c[j].SkuTransacted()) },
	)

	return
}

func (s *zettelStore) addZettelToTransaktion(
	z zettel.Named,
) (tz zettel.Transacted, err error) {
	errors.Log().Printf("adding zettel to transaktion: %s", z.Kennung)

	if tz, err = s.transactedWithHead(z, &s.common.Transaktion); err != nil {
		err = errors.Wrap(err)
		return
	}

	i := s.common.Transaktion.Add(
		sku.Sku{
			Gattung: gattung.Zettel,
			Mutter:  tz.Mutter,
			Id:      &z.Kennung,
			Sha:     z.Stored.Sha,
		},
	)

	tz.TransaktionIndex.SetInt(i)

	return

}
func (s zettelStore) storedZettelFromSha(sh sha.Sha) (sz zettel.Stored, err error) {
	var or io.ReadCloser

	if or, err = s.common.ReadCloserObjekten(id.Path(sh, s.common.Standort.DirObjektenZettelen())); err != nil {
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
func (s *zettelStore) transactedWithHead(
	z zettel.Named,
	t *transaktion.Transaktion,
) (tz zettel.Transacted, err error) {
	tz.Named = z
	tz.Kopf = t.Time
	tz.Schwanz = t.Time

	var previous zettel.Transacted

	if previous, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(z.Kennung); err == nil {
		tz.Mutter[0] = previous.Schwanz
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

func (s zettelStore) transactedZettelFromTransaktionObjekte(
	t *transaktion.Transaktion,
	o *sku.Indexed,
) (tz zettel.Transacted, err error) {
	ok := false

	var h *hinweis.Hinweis

	if h, ok = o.Id.(*hinweis.Hinweis); !ok {
		err = errors.Wrapf(err, "transaktion.Objekte Id was not hinweis but was %s", o.Id)
		return
	}

	tz.Named.Kennung = *h

	if tz.Named.Stored, err = s.storedZettelFromSha(o.Sha); err != nil {
		err = errors.Wrapf(err, "failed to read zettel objekte: %s", tz.Named.Kennung)
		return
	}

	if tz, err = s.transactedWithHead(tz.Named, t); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz.TransaktionIndex = o.Index

	return
}
