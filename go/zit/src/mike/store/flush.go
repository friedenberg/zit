package store

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/bestandsaufnahme"
)

func (s *Store) FlushBestandsaufnahme(
	p schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	if s.GetKonfig().DryRun {
		return
	}

	if !s.GetStandort().GetLockSmith().IsAcquired() {
		return
	}

	ui.Log().Printf("saving Bestandsaufnahme")

	var bestandsaufnahmeSku *sku.Transacted

	if bestandsaufnahmeSku, err = s.GetBestandsaufnahmeStore().Create(
		&s.bestandsaufnahmeAkte,
		s.GetKonfig().Bezeichnung,
	); err != nil {
		if errors.Is(err, bestandsaufnahme.ErrEmpty) {
			ui.Log().Printf("Bestandsaufnahme was empty")
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if bestandsaufnahmeSku != nil {
		defer sku.GetTransactedPool().Put(bestandsaufnahmeSku)

		if s.GetKonfig().PrintOptions.PrintBestandsaufnahme {
			if err = p(bestandsaufnahmeSku); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	bestandsaufnahme.Resetter.Reset(&s.bestandsaufnahmeAkte)

	if err = s.GetBestandsaufnahmeStore().Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ui.Log().Printf("done saving Bestandsaufnahme")

	return
}

func (c *Store) Flush(
	printerHeader schnittstellen.FuncIter[string],
) (err error) {
	// TODO handle flushes with dry run
	if c.GetKonfig().DryRun {
		return
	}

	wg := iter.MakeErrorWaitGroupParallel()

	if c.GetStandort().GetLockSmith().IsAcquired() {
		gob.Register(iter.StringerKeyerPtr[kennung.Typ, *kennung.Typ]{}) // TODO check if can be removed
		wg.Do(func() error { return c.verzeichnisse.Flush(printerHeader) })
		wg.Do(c.GetAbbrStore().Flush)
		wg.Do(c.typenIndex.Flush)
		wg.Do(c.kennungIndex.Flush)
		wg.Do(c.Abbr.Flush)
	}

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
