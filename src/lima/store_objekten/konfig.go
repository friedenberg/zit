package store_objekten

import (
	"encoding/gob"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/echo/age_io"
	"github.com/friedenberg/zit/src/echo/konfig"
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

func (s konfigStore) transact(
	tt *konfig.Transacted,
	mutterKopf ts.Time,
	mutterSchanz ts.Time,
) (err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "write named zettel to index",
		}

		return
	}

	tt.Schwanz = s.common.Transaktion.Time
	tt.Kopf = mutterKopf
	tt.Mutter[0] = mutterSchanz

	sk := tt.SkuTransacted()

	tt.TransaktionIndex.SetInt(s.common.Transaktion.Add(sk.Sku))

	var f *os.File

	if f, err = files.OpenExclusiveWriteOnly(s.common.Standort.FileKonfigCompiled()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	enc := gob.NewEncoder(f)

	if err = enc.Encode(tt); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

//This is intentionally a value to prevent accidental global changes to the
//current Konfig
func (s konfigStore) Read() (tt konfig.Transacted, err error) {
	tt = s.common.Konfig.Transacted

	return
}

// TODO support dry run
func (s *konfigStore) Update(
	t *konfig.Objekte,
) (tt *konfig.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "update",
		}

		return
	}

	var mutter konfig.Transacted

	if mutter, err = s.Read(); err != nil {
		err = errors.Wrap(err)
		return
	}

	tt = &konfig.Transacted{
		Named: konfig.Named{
			Stored: konfig.Stored{
				Objekte: *t,
			},
		},
	}

	if err = s.writeObjekte(&tt.Named.Stored); err != nil {
		err = errors.Wrap(err)
		return
	}

	if tt.Named.Stored.Sha.Equals(mutter.Named.Stored.Sha) {
		tt = &mutter

		if err = s.KonfigLogWriters.Unchanged(tt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = s.transact(
		tt,
		mutter.Kopf,
		mutter.Schwanz,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.KonfigLogWriters.Updated(tt); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s konfigStore) AllInChain() (c []*konfig.Transacted, err error) {

	return
}
