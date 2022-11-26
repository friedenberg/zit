package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/age_io"
	"github.com/friedenberg/zit/src/golf/typ"
	"github.com/friedenberg/zit/src/typ_toml"
)

type typLogWriter = collections.WriterFunc[*typ.Transacted]

type TypLogWriters struct {
	New, Updated, Archived, Unchanged typLogWriter
}

type typStore struct {
	common *common
	TypLogWriters
}

func (s *typStore) SetTypLogWriters(
	tlw TypLogWriters,
) {
	s.TypLogWriters = tlw
}

func makeTypStore(
	sa *common,
) (s *typStore, err error) {
	s = &typStore{
		common: sa,
	}

	return
}

func (s typStore) Flush() (err error) {
	return
}

func (s typStore) writeObjekte(t *typ.Stored) (sh sha.Sha, err error) {
	//no lock required

	var w *age_io.Mover

	mo := age_io.MoveOptions{
		Age:                      s.common.Age,
		FinalPath:                s.common.Standort.DirObjektenTypen(),
		GenerateFinalPathFromSha: true,
	}

	if w, err = age_io.NewMover(mo); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer w.Close()

	f := typ.MakeFormatObjekte(s.common)

	if _, err = f.WriteFormat(w, t); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = w.Sha()

	return
}

//TODO write konfig compiled
func (s typStore) writeTransactedToIndex(tt *typ.Transacted) (err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "write named zettel to index",
		}

		return
	}

	return
}

//TODO
func (s typStore) ReadOne(
	k kennung.Typ,
) (tt *typ.Transacted, err error) {
	ct := s.common.Konfig.Compiled.GetTyp(k.String())

	if ct == nil {
		err = errors.Wrap(ErrNotFound{Id: k})
		return
	}

	tt = &typ.Transacted{
		Named: typ.Named{
			Kennung: k,
			Stored: typ.Stored{
				//TODO
				// Sha: sha,
				Objekte: typ.Akte{
					KonfigTyp: typ_toml.Typ{
						InlineAkte:     ct.InlineAkte,
						FileExtension:  ct.FileExtension,
						ExecCommand:    ct.ExecCommand,
						Actions:        ct.Actions,
						EtikettenRules: ct.EtikettenRules,
					},
				},
			},
		},
	}

	return
}

func (s *typStore) Create(in typ.Akte, k kennung.Typ) (tt *typ.Transacted, err error) {
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

	tt = &typ.Transacted{
		Named: typ.Named{
			Stored: typ.Stored{
				Objekte: in,
			},
			Kennung: k,
		},
	}

	if tt.Named.Stored.Sha, err = s.writeObjekte(&tt.Named.Stored); err != nil {
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
	if err = s.TypLogWriters.New(tt); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO support dry run
func (s *typStore) Update(
	t *typ.Named,
) (tt *typ.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "update",
		}

		return
	}

	var mutter *typ.Transacted

	if mutter, err = s.ReadOne(
		t.Kennung,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if t.Equals(&mutter.Named) {
		tt = mutter

		if err = s.TypLogWriters.Unchanged(tt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	tt.Named = *t

	if tt.Named.Stored.Sha, err = s.writeObjekte(&t.Stored); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.writeTransactedToIndex(tt); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.TypLogWriters.Updated(tt); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s typStore) AllInChain(k kennung.Typ) (c []*typ.Transacted, err error) {
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
