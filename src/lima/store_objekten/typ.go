package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/foxtrot/typ_toml"
	"github.com/friedenberg/zit/src/golf/age_io"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/india/typ"
)

type typLogWriter = collections.WriterFunc[*typ.Transacted]

type TypLogWriters struct {
	New, Updated, Unchanged typLogWriter
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

func (s typStore) transact(
	to *typ.Objekte,
	tk *kennung.Typ,
) (tt *typ.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "transact typ",
		}

		return
	}

	var mutter *typ.Transacted

	if mutter, err = s.ReadOne(tk); err != nil {
		if errors.Is(err, ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	tt = &typ.Transacted{
		Objekte: *to,
		Sku: sku.Transacted[kennung.Typ, *kennung.Typ]{
			Kennung: *tk,
			Schwanz: s.common.Transaktion.Time,
		},
	}

	//TODO-P3 refactor into reusable
	if mutter != nil {
		tt.Sku.Kopf = mutter.Sku.Kopf
		tt.Sku.Mutter[0] = mutter.Sku.Schwanz
	} else {
		tt.Sku.Kopf = s.common.Transaktion.Time
	}

	fo := objekte.MakeFormatter[*typ.Transacted](s.common)

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

	defer errors.Deferred(&err, w.Close)

	if _, err = fo.WriteFormat(w, tt); err != nil {
		err = errors.Wrap(err)
		return
	}

	tt.Sku.Sha = w.Sha()

	if mutter != nil && tt.ObjekteSha().Equals(mutter.ObjekteSha()) {
		tt = mutter

		if err = s.TypLogWriters.Unchanged(tt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.common.Transaktion.Add2(&tt.Sku)

	if mutter == nil {
		if err = s.TypLogWriters.New(tt); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = s.TypLogWriters.Updated(tt); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

// TODO-P0 disambiguate from akte
func (s typStore) WriteAkte(
	t *typ.Objekte,
) (err error) {
	var w sha.WriteCloser

	if w, err = s.common.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, w.Close)

	if _, err = typ_toml.WriteObjekteToText(w, t); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.Sha = w.Sha()

	return
}

func (s typStore) ReadOne(
	k *kennung.Typ,
) (tt *typ.Transacted, err error) {
	tt = s.common.Konfig.GetTyp(k.String())

	if tt == nil {
		err = errors.Wrap(ErrNotFound{Id: k})
		return
	}

	return
}

func (s *typStore) Create(
	in *typ.Objekte,
	tk *kennung.Typ,
) (tt *typ.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "create typ",
		}

		return
	}

	if tt, err = s.transact(in, tk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *typStore) Update(
	t *typ.Objekte,
	tk *kennung.Typ,
) (tt *typ.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "update typ",
		}

		return
	}

	if tt, err = s.transact(t, tk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s typStore) AllInChain(k kennung.Typ) (c []*typ.Transacted, err error) {
	// mst := zettel_transacted.MakeMutableSetUnique(0)

	// if err = s.verzeichnisseAll.ReadMany(
	// 	func(z *zettel_verzeichnisse.Zettel) (err error) {
	// 		if !z.Transacted.Sku.Kennung.Equals(&h) {
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
