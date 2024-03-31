package store

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/juliett/objekte"
	"code.linenisgreat.com/zit/src/lima/bestandsaufnahme"
)

func (s *Store) FlushBestandsaufnahme() (err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte.ErrLockRequired{
			Operation: "flush",
		}

		return
	}

	if s.GetKonfig().DryRun {
		return
	}

	errors.Log().Printf("saving Bestandsaufnahme")

	if _, err = s.GetBestandsaufnahmeStore().Create(
		&s.bestandsaufnahmeAkte,
	); err != nil {
		if errors.Is(err, bestandsaufnahme.ErrEmpty) {
			errors.Log().Printf("Bestandsaufnahme was empty")
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	bestandsaufnahme.Resetter.Reset(&s.bestandsaufnahmeAkte)

	if err = s.GetBestandsaufnahmeStore().Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("done saving Bestandsaufnahme")

	return
}

func (c *Store) Flush(
	printerHeader schnittstellen.FuncIter[string],
) (err error) {
	if !c.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte.ErrLockRequired{
			Operation: "flush",
		}

		return
	}

	if c.GetKonfig().DryRun {
		return
	}

	gob.Register(iter.StringerKeyerPtr[kennung.Typ, *kennung.Typ]{})

	if c.GetKonfig().HasChanges() {
		c.verzeichnisse.SetNeedsFlushHistory()
	}

	wg := iter.MakeErrorWaitGroupParallel()

	wg.Do(func() error { return c.verzeichnisse.Flush(printerHeader) })
	wg.Do(c.GetAbbrStore().Flush)
	wg.Do(c.typenIndex.Flush)
	wg.Do(c.kennungIndex.Flush)
	wg.Do(c.Abbr.Flush)

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
