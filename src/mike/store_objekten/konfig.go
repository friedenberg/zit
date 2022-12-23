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
	"github.com/friedenberg/zit/src/india/konfig"
)

type KonfigLogWriter = collections.WriterFunc[*konfig.Transacted]

type KonfigLogWriters struct {
	Updated, Unchanged KonfigLogWriter
}

type konfigStore struct {
	common *common

	objekte.TransactedInflator[
		konfig.Objekte,
		*konfig.Objekte,
		kennung.Konfig,
		*kennung.Konfig,
	]

	objekte.AkteTextSaver[
		konfig.Objekte,
		*konfig.Objekte,
	]

	KonfigLogWriters
}

func (s *konfigStore) SetKonfigLogWriters(
	tlw KonfigLogWriters,
) {
	s.KonfigLogWriters = tlw
}

func makeKonfigStore(
	sa *common,
) (s *konfigStore, err error) {
	s = &konfigStore{
		common: sa,
		TransactedInflator: objekte.MakeTransactedInflator[
			konfig.Objekte,
			*konfig.Objekte,
			kennung.Konfig,
			*kennung.Konfig,
		](
			sa,
			func(sh sha.Sha) (r sha.ReadCloser, err error) {
				return s.common.ReadCloserObjekten(
					id.Path(sh, sa.Standort.DirObjektenKonfig()),
				)
			},
			nil,
			gattung.Parser[konfig.Objekte, *konfig.Objekte](
				konfig.MakeFormatText(sa),
			),
		),
		AkteTextSaver: objekte.MakeAkteTextSaver[
			konfig.Objekte,
			*konfig.Objekte,
		](
			sa,
			&konfig.FormatterAkteTextToml{},
		),
	}

	return
}

func (s konfigStore) Flush() (err error) {
	return
}

func (s konfigStore) Update(
	ko *konfig.Objekte,
) (kt *konfig.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = errors.Wrap(ErrLockRequired{Operation: "update konfig"})
		return
	}

	var w *age_io.Mover

	mo := age_io.MoveOptions{
		Age:                      s.common.Age,
		FinalPath:                s.common.Standort.DirObjektenKonfig(),
		GenerateFinalPathFromSha: true,
	}

	if w, err = age_io.NewMover(mo); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, w.Close)

	var mutter *konfig.Transacted

	if mutter, err = s.Read(); err != nil {
		if errors.Is(err, ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	kt = &konfig.Transacted{
		Objekte: *ko,
		Sku: sku.Transacted[kennung.Konfig, *kennung.Konfig]{
			Schwanz: s.common.Transaktion.Time,
		},
	}

	//TODO-P3 refactor into reusable
	if mutter != nil {
		kt.Sku.Kopf = mutter.Sku.Kopf
		kt.Sku.Mutter[0] = mutter.Sku.Schwanz
	} else {
		kt.Sku.Kopf = s.common.Transaktion.Time
	}

	fo := objekte.MakeFormat[konfig.Objekte, *konfig.Objekte]()

	if _, err = fo.Format(w, &kt.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	kt.Sku.Sha = w.Sha()

	if mutter != nil && kt.ObjekteSha().Equals(mutter.ObjekteSha()) {
		kt = mutter

		if err = s.KonfigLogWriters.Unchanged(kt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.common.Transaktion.Add2(&kt.Sku)
	s.common.KonfigPtr().SetTransacted(kt)

	if err = s.common.Abbr.addStored(kt); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.KonfigLogWriters.Updated(kt); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s konfigStore) Read() (tt *konfig.Transacted, err error) {
	tt = &konfig.Transacted{
		Sku: s.common.Konfig().Sku,
		Objekte: konfig.Objekte{
			Akte: s.common.Konfig().Toml,
		},
	}

	if !tt.Sku.Schwanz.IsEmpty() {
		{
			var r sha.ReadCloser

			if r, err = s.common.ReadCloserObjekten(
				id.Path(tt.Sku.Sha, s.common.Standort.DirObjektenKonfig()),
			); err != nil {
				if errors.IsNotExist(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
				}

				return
			}

			defer errors.Deferred(&err, r.Close)

			fo := objekte.MakeFormat[konfig.Objekte, *konfig.Objekte]()

			if _, err = fo.Parse(r, &tt.Objekte); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		{
			var r sha.ReadCloser

			if r, err = s.common.ReadCloserObjekten(
				id.Path(tt.Objekte.Sha, s.common.Standort.DirObjektenAkten()),
			); err != nil {
				if errors.IsNotExist(err) {
					err = errors.Wrap(ErrNotFound{})
				} else {
					err = errors.Wrap(err)
				}
				return
			}

			defer errors.Deferred(&err, r.Close)

			fo := konfig.MakeFormatText(s.common)

			if _, err = fo.ReadFormat(r, &tt.Objekte); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

func (s konfigStore) AllInChain() (c []*konfig.Transacted, err error) {
	return
}

func (s *konfigStore) reindexOne(
	t *transaktion.Transaktion,
	o *sku.Sku,
) (err error) {
	var te *konfig.Transacted

	if te, err = s.Inflate(t, o); err != nil {
		errors.Wrap(err)
		return
	}

	s.common.KonfigPtr().SetTransacted(te)
	s.KonfigLogWriters.Updated(te)

	return
}
