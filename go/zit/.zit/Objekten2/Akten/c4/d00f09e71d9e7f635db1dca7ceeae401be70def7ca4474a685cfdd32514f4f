package store

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/inventory_list_store"
)

func (s *Store) FlushInventoryList(
	p interfaces.FuncIter[*sku.Transacted],
) (err error) {
	if s.GetConfig().DryRun {
		return
	}

	if !s.GetDirectoryLayout().GetLockSmith().IsAcquired() {
		return
	}

	ui.Log().Printf("saving Bestandsaufnahme")

	var inventoryListSku *sku.Transacted

	if inventoryListSku, err = s.GetInventoryListStore().Create(
		s.inventoryList,
		s.GetConfig().Description,
	); err != nil {
		if errors.Is(err, inventory_list_store.ErrEmpty) {
			ui.Log().Printf("Bestandsaufnahme was empty")
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if inventoryListSku != nil {
		if err = s.GetStreamIndex().Add(
			inventoryListSku,
			inventoryListSku.GetObjectId().String(),
			sku.CommitOptions{
				Mode: object_mode.ModeLatest,
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
		defer sku.GetTransactedPool().Put(inventoryListSku)

		if s.GetConfig().PrintOptions.PrintInventoryLists {
			if err = p(inventoryListSku); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	sku.ResetterList.Reset(s.inventoryList)

	if err = s.GetInventoryListStore().Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ui.Log().Printf("done saving Bestandsaufnahme")

	return
}

func (c *Store) Flush(
	printerHeader interfaces.FuncIter[string],
) (err error) {
	// TODO handle flushes with dry run
	if c.GetConfig().DryRun {
		return
	}

	wg := quiter.MakeErrorWaitGroupParallel()

	if c.GetDirectoryLayout().GetLockSmith().IsAcquired() {
		gob.Register(quiter.StringerKeyerPtr[ids.Type, *ids.Type]{}) // TODO check if can be removed
		wg.Do(func() error { return c.streamIndex.Flush(printerHeader) })
		wg.Do(c.GetAbbrStore().Flush)
		wg.Do(c.zettelIdIndex.Flush)
		wg.Do(c.Abbr.Flush)
	}

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
