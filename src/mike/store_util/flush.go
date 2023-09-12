package store_util

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/kilo/bestandsaufnahme"
	"github.com/friedenberg/zit/src/lima/objekte_store"
)

func (s *common) FlushBestandsaufnahme() (err error) {
	if !s.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "flush",
		}

		return
	}

	if s.GetKonfig().DryRun {
		return
	}

	errors.Log().Printf("saving Bestandsaufnahme")
	ba := s.GetBestandsaufnahmeAkte()
	if err = s.GetBestandsaufnahmeStore().Create(&ba); err != nil {
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
