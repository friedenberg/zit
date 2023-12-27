package store_util

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/echo/kennung"
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

	if _, err = s.GetBestandsaufnahmeStore().Create(&s.bestandsaufnahmeAkte); err != nil {
		if errors.Is(err, bestandsaufnahme.ErrEmpty) {
			errors.Log().Printf("Bestandsaufnahme was empty")
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	bestandsaufnahme.Resetter.Reset(&s.bestandsaufnahmeAkte)

	errors.Log().Printf("done saving Bestandsaufnahme")

	return
}

func (c *common) Flush() (err error) {
	if c.konfig.DryRun {
		return
	}

	gob.Register(iter.StringerKeyerPtr[kennung.Typ, *kennung.Typ]{})

	if c.GetKonfig().HasChanges() {
		c.verzeichnisse.SetNeedsFlush()
	}

	wg := iter.MakeErrorWaitGroupParallel()

	wg.Do(c.verzeichnisse.Flush)
	wg.Do(c.typenIndex.Flush)
	wg.Do(c.kennungIndex.Flush)
	wg.Do(c.Abbr.Flush)

	err = wg.GetError()

	return
}
