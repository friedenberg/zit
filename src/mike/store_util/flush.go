package store_util

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
	"github.com/friedenberg/zit/src/lima/bestandsaufnahme"
)

func (s *common) FlushBestandsaufnahme() (err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "flush",
		}

		return
	}

	if s.GetKonfig().DryRun {
		return
	}

	errors.Log().Printf("saving Bestandsaufnahme")
	if err = s.GetBestandsaufnahmeStore().Create(&s.bestandsaufnahmeAkte); err != nil {
		if errors.Is(err, bestandsaufnahme.ErrEmpty) {
			errors.Log().Printf("Bestandsaufnahme was empty")
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	errors.Log().Printf("done saving Bestandsaufnahme")

	return
}

func (c *common) Flush() (err error) {
	if err = c.typenIndex.Flush(); err != nil {
		err = errors.Wrapf(err, "failed to flush typen index")
		return
	}

	if err = c.kennungIndex.Flush(); err != nil {
		err = errors.Wrapf(err, "failed to flush kennung index")
		return
	}

	if err = c.Abbr.Flush(); err != nil {
		err = errors.Wrapf(err, "failed to flush abbr index")
		return
	}

	return
}
