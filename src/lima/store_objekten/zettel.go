package store_objekten

import (
	"io"
	"reflect"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/hinweisen"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/echo/age_io"
	"github.com/friedenberg/zit/src/india/zettel"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/zettel_verzeichnisse"
)

// TODO add archived state
type ZettelTransactedLogWriters struct {
	New, Updated, Archived, Unchanged collections.WriterFunc[*zettel_transacted.Zettel]
}

func (s *Store) SetZettelTransactedLogWriter(
	ztlw ZettelTransactedLogWriters,
) {
	s.zettelTransactedWriter = ztlw
}

func (s Store) WriteZettelObjekte(z zettel.Zettel) (sh sha.Sha, err error) {
	//no lock required

	var w *age_io.Mover

	mo := age_io.MoveOptions{
		Age:                      s.age,
		FinalPath:                s.standort.DirObjektenZettelen(),
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

func (s Store) writeNamedZettelToIndex(tz zettel_transacted.Zettel) (err error) {
	if !s.lockSmith.IsAcquired() {
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

func (s Store) ReadHinweisSchwanzen(
	h hinweis.Hinweis,
) (zv zettel_transacted.Zettel, err error) {
	return s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(h)
}

func (i *Store) ReadAllSchwanzenVerzeichnisse(
	ws ...collections.WriterFunc[*zettel_verzeichnisse.Zettel],
) (err error) {
	return i.verzeichnisseSchwanzen.ReadMany(ws...)
}

func (s Store) ReadAllSchwanzenTransacted(
	ws ...collections.WriterFunc[*zettel_transacted.Zettel],
) (err error) {
	w := zettel_verzeichnisse.MakeWriterZettelTransacted(
		collections.MakeChain(ws...),
	)

	return s.ReadAllSchwanzenVerzeichnisse(w)
}

func (i *Store) ReadAllVerzeichnisse(
	ws ...collections.WriterFunc[*zettel_verzeichnisse.Zettel],
) (err error) {
	return i.verzeichnisseAll.ReadMany(ws...)
}

func (s Store) ReadAllTransacted(
	ws ...collections.WriterFunc[*zettel_transacted.Zettel],
) (err error) {
	w := zettel_verzeichnisse.MakeWriterZettelTransacted(
		collections.MakeChain(ws...),
	)

	return s.ReadAllVerzeichnisse(w)
}

func (s Store) ReadOne(i id.Id) (tz zettel_transacted.Zettel, err error) {
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

func (s *Store) Create(in zettel.Zettel) (tz zettel_transacted.Zettel, err error) {
	if !s.lockSmith.IsAcquired() {
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

	if err = in.ApplyKonfig(s.konfig); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz.Named.Stored.Objekte = in

	if tz.Named.Stored.Sha, err = s.WriteZettelObjekte(tz.Named.Stored.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO-P1?
	//If the zettel exists, short circuit and return that
	// if tz2, err2 := s.Read(tz.Named.Stored.Sha); err2 == nil {
	// 	tz = tz2
	// 	return
	// }

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

func (s *Store) CreateWithHinweis(
	in zettel.Zettel,
	h hinweis.Hinweis,
) (tz zettel_transacted.Zettel, err error) {
	if !s.lockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "create with hinweis",
		}

		return
	}

	if in.IsEmpty() {
		err = errors.Normalf("zettel is empty")
		return
	}

	if err = in.ApplyKonfig(s.konfig); err != nil {
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
func (s *Store) Update(
	z *zettel.Named,
) (tz zettel_transacted.Zettel, err error) {
	if !s.lockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "update",
		}

		return
	}

	if err = z.Stored.Objekte.ApplyKonfig(s.konfig); err != nil {
		err = errors.Wrap(err)
		return
	}

	var mutter zettel_transacted.Zettel

	if mutter, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(
		z.Kennung,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO-P1 prevent useless duplicate transaktions
	// if z.Equals(mutter.Named) {
	// 	tz = mutter
	// 	return
	// }

	tz.Named = *z

	if tz.Named.Stored.Sha, err = s.WriteZettelObjekte(z.Stored.Objekte); err != nil {
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

	if err = s.indexEtiketten.addZettelWithOptionalMutter(&tz, &mutter); err != nil {
		err = errors.Wrap(err)
		return
	}

	if mutter.Named.Equals(&tz.Named) {
		if err = s.zettelTransactedWriter.Unchanged(&tz); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = s.zettelTransactedWriter.Updated(&tz); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s Store) AllInChain(h hinweis.Hinweis) (c []*zettel_transacted.Zettel, err error) {
	mst := zettel_transacted.MakeMutableSetUnique(0)

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
