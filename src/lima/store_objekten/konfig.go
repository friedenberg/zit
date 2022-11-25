package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/age_io"
	"github.com/friedenberg/zit/src/echo/konfig"
)

type konfigLogWriter = collections.WriterFunc[*konfig.Transacted]

type konfigLogWriters struct {
	New, Updated, Archived, Unchanged konfigLogWriter
}

type konfigStore struct {
	common *common
	konfigLogWriters
}

func (s *konfigStore) SetkonfigLogWriters(
	tlw konfigLogWriters,
) {
	s.konfigLogWriters = tlw
}

func makekonfigStore(
	sa *common,
) (s *konfigStore, err error) {
	s = &konfigStore{
		common: sa,
	}

	return
}

func (s konfigStore) Flush() (err error) {
	return
}

func (s konfigStore) writeObjekte(t *konfig.Stored) (err error) {
	//no lock required

	var w *age_io.Mover

	mo := age_io.MoveOptions{
		Age:                      s.common.Age,
		FinalPath:                s.common.Standort.DirObjektenAkten(),
		GenerateFinalPathFromSha: true,
	}

	if w, err = age_io.NewMover(mo); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer w.Close()

	f := konfig.MakeFormatObjekte(s.common)

	if _, err = f.WriteFormat(w, t); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.Sha = w.Sha()

	return
}

//TODO write konfig compiled
func (s konfigStore) writeTransactedToIndex(tt *konfig.Transacted) (err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "write named zettel to index",
		}

		return
	}

	return
}

//TODO
func (s konfigStore) ReadOne(
	k kennung.Konfig,
) (tt *konfig.Transacted, err error) {
	ct := s.common.Konfig.GetTyp(k.String())

	if ct == nil {
		err = errors.Wrap(ErrNotFound{Id: k})
		return
	}

	tt = &konfig.Transacted{
		Named: konfig.Named{
			Kennung: k,
			Stored: konfig.Stored{
				//TODO
				// Sha: sha,
				Objekte: konfig.Konfig{},
			},
		},
	}

	return
}

func (s *konfigStore) Create(in konfig.Konfig) (tt *konfig.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "create",
		}

		return
	}

	//TODO
	// if in.IsEmpty() {
	// 	err = errors.Normalf("%s is empty", in.Gattung())
	// 	return
	// }

	tt = &konfig.Transacted{
		Named: konfig.Named{
			Stored: konfig.Stored{
				Objekte: in,
			},
			Kennung: kennung.Konfig{},
		},
	}

	if err = s.writeObjekte(&tt.Named.Stored); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO-P1?
	//If the zettel exists, short circuit and return that
	// if tz2, err2 := s.Read(tz.Named.Stored.Sha); err2 == nil {
	// 	tz = tz2
	// 	return
	// }

	//TODO?
	// if tz, err = s.addZettelToTransaktion(tz.Named); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	if err = s.writeTransactedToIndex(tt); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO-P2 assert no changes
	if err = s.konfigLogWriters.New(tt); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO support dry run
func (s *konfigStore) Update(
	t *konfig.Named,
) (tt *konfig.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "update",
		}

		return
	}

	var mutter *konfig.Transacted

	if mutter, err = s.ReadOne(
		t.Kennung,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if t.Equals(&mutter.Named) {
		tt = mutter

		if err = s.konfigLogWriters.Unchanged(tt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	tt.Named = *t

	if err = s.writeObjekte(&t.Stored); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.writeTransactedToIndex(tt); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.konfigLogWriters.Updated(tt); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s konfigStore) AllInChain() (c []*konfig.Transacted, err error) {
	// mst := zettel_transacted.MakeMutableSetUnique(0)

	// if err = s.verzeichnisseAll.ReadMany(
	// 	func(z *zettel_verzeichnisse.Zettel) (err error) {
	// 		if !z.Transacted.Named.Kennung.Equals(&h) {
	// 			err = io.EOF
	// 			return
	// 		}

	// 		return
	// 	},
	// 	zettel_verzeichnisse.MakeWriterZettelTransacted(mst.AddAndDoNotRepool),
	// ); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	// c = mst.Elements()

	// sort.Slice(
	// 	c,
	// 	func(i, j int) bool { return c[i].SkuTransacted().Less(c[j].SkuTransacted()) },
	// )

	return
}
