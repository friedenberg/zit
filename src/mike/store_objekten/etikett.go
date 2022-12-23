package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/id"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/golf/age_io"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/india/etikett"
)

type etikettLogWriter = collections.WriterFunc[*etikett.Transacted]

type EtikettLogWriters struct {
	New, Updated, Unchanged etikettLogWriter
}

type etikettStore struct {
	common *common

	objekte.Inflator[
		etikett.Objekte,
		*etikett.Objekte,
		kennung.Etikett,
		*kennung.Etikett,
	]

	EtikettLogWriters
}

func (s *etikettStore) SetEtikettLogWriters(
	tlw EtikettLogWriters,
) {
	s.EtikettLogWriters = tlw
}

func makeEtikettStore(
	sa *common,
) (s *etikettStore, err error) {
	s = &etikettStore{
		common: sa,
		Inflator: objekte.MakeTransactedInflator[
			etikett.Objekte,
			*etikett.Objekte,
			kennung.Etikett,
			*kennung.Etikett,
		](
			sa,
			func(sh sha.Sha) (r sha.ReadCloser, err error) {
				return s.common.ReadCloserObjekten(
					id.Path(sh, sa.Standort.DirObjektenEtiketten()),
				)
			},
			gattung.Parser[etikett.Objekte, *etikett.Objekte](
				etikett.MakeFormatText(sa),
			),
		),
	}

	return
}

func (s etikettStore) Flush() (err error) {
	return
}

func (s etikettStore) transact(
	to *etikett.Objekte,
	tk *kennung.Etikett,
) (tt *etikett.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "transact etikett",
		}

		return
	}

	var mutter *etikett.Transacted

	if mutter, err = s.ReadOne(tk); err != nil {
		if errors.Is(err, ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	tt = &etikett.Transacted{
		Objekte: *to,
		Sku: sku.Transacted[kennung.Etikett, *kennung.Etikett]{
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

	fo := objekte.MakeFormatter[*etikett.Transacted](s.common)

	var w *age_io.Mover

	mo := age_io.MoveOptions{
		Age:                      s.common.Age,
		FinalPath:                s.common.Standort.DirObjektenEtiketten(),
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

		if err = s.EtikettLogWriters.Unchanged(tt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.common.Transaktion.Add2(&tt.Sku)
	s.common.KonfigPtr().AddEtikett(tt)

	if mutter == nil {
		if err = s.EtikettLogWriters.New(tt); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = s.EtikettLogWriters.Updated(tt); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

// TODO-P0 disambiguate from akte
func (s etikettStore) WriteAkte(
	t *etikett.Objekte,
) (err error) {
	var w sha.WriteCloser

	if w, err = s.common.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, w.Close)

	if _, err = etikett.WriteObjekteToText(w, t); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.Sha = w.Sha()

	return
}

func (s etikettStore) ReadOne(
	k *kennung.Etikett,
) (tt *etikett.Transacted, err error) {
	tt = s.common.Konfig().GetEtikett(*k)

	if tt == nil {
		err = errors.Wrap(ErrNotFound{Id: k})
		return
	}

	return
}

func (s *etikettStore) Create(
	in *etikett.Objekte,
	tk *kennung.Etikett,
) (tt *etikett.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "create etikett",
		}

		return
	}

	if tt, err = s.transact(in, tk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *etikettStore) Update(
	t *etikett.Objekte,
	tk *kennung.Etikett,
) (tt *etikett.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "update etikett",
		}

		return
	}

	if tt, err = s.transact(t, tk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s etikettStore) AllInChain(k kennung.Etikett) (c []*etikett.Transacted, err error) {
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

func (s *etikettStore) reindexOne(
	t *transaktion.Transaktion,
	o *sku.Sku,
) (err error) {
	var te *etikett.Transacted

	if te, err = s.Inflate(t, o); err != nil {
		errors.Wrap(err)
		return
	}

	s.common.KonfigPtr().AddEtikett(te)

	if te.IsNew() {
		s.EtikettLogWriters.New(te)
	} else {
		s.EtikettLogWriters.Updated(te)
	}

	return
}
