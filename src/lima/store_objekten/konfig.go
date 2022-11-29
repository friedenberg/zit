package store_objekten

import (
	"encoding/gob"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/echo/age_io"
	"github.com/friedenberg/zit/src/echo/konfig"
	"github.com/friedenberg/zit/src/echo/sku"
	"github.com/friedenberg/zit/src/foxtrot/objekte"
)

type KonfigLogWriter = collections.WriterFunc[*konfig.Transacted]

type KonfigLogWriters struct {
	Updated, Unchanged KonfigLogWriter
}

type konfigStore struct {
	common *common
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
	}

	return
}

func (s konfigStore) Flush() (err error) {
	return
}

func (s konfigStore) transact(
	ko *konfig.Objekte,
) (kt *konfig.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "transact konfig",
		}

		return
	}

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

	var mutter konfig.Transacted

	if mutter, err = s.Read(); err != nil {
		err = errors.Wrap(err)
		return
	}

	kt = &konfig.Transacted{
		Objekte: *ko,
		Sku: sku.Sku2[kennung.Konfig, *kennung.Konfig]{
			Schwanz: s.common.Transaktion.Time,
			Kopf:    mutter.Sku.Kopf,
			Mutter:  sku.Mutter{mutter.Sku.Schwanz, ts.Time{}},
		},
	}

	fo := objekte.MakeFormatObjekte(s.common)

	if _, err = fo.WriteFormat(w, kt); err != nil {
		err = errors.Wrap(err)
		return
	}

	kt.Sku.Sha = w.Sha()

	if kt.ObjekteSha().Equals(mutter.ObjekteSha()) {
		kt = &mutter

		if err = s.KonfigLogWriters.Unchanged(kt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.common.Transaktion.Add2(&kt.Sku)

	if !s.common.Konfig.DryRun {
		var f *os.File

		//TODO use objekte mover
		if f, err = files.OpenExclusiveWriteOnly(s.common.Standort.FileKonfigCompiled()); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, f.Close)

		enc := gob.NewEncoder(f)

		if err = enc.Encode(kt); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.KonfigLogWriters.Updated(kt); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// This is intentionally a value to prevent accidental global changes to the
// current Konfig
func (s konfigStore) Read() (tt konfig.Transacted, err error) {
	tt = s.common.Konfig.Transacted

	return
}

func (s *konfigStore) Update(
	ko *konfig.Objekte,
) (kt *konfig.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "update",
		}

		return
	}

	if kt, err = s.transact(ko); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s konfigStore) AllInChain() (c []*konfig.Transacted, err error) {

	return
}
