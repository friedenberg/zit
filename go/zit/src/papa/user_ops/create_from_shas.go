package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type CreateFromShas struct {
	*env.Local
	sku.Proto
}

func (c CreateFromShas) Run(
	args ...string,
) (results sku.TransactedMutableSet, err error) {
	var lookupStored map[sha.Bytes][]string

	if lookupStored, err = c.GetStore().MakeBlobShaBytesMap(); err != nil {
		err = errors.Wrap(err)
		return
	}

	toCreate := make(map[sha.Bytes]*sku.Transacted)

	for _, arg := range args {
		var sh sha.Sha

		if err = sh.Set(arg); err != nil {
			err = errors.Wrap(err)
			return
		}

		k := sh.GetBytes()

		if _, ok := toCreate[k]; ok {
			ui.Err().Printf("%s appears in arguments more than once. Ignoring", &sh)
			continue
		}

		if oids, ok := lookupStored[k]; ok {
			ui.Err().Printf("%s appears in object already checked in (%q). Ignoring", &sh, oids)
			continue
		}

		z := sku.GetTransactedPool().Get()

		z.ObjectId.SetGenre(genres.Zettel)
		z.Metadata.Blob.ResetWith(&sh)

		c.Proto.Apply(z, genres.Zettel)

		toCreate[k] = z
	}

	results = sku.MakeTransactedMutableSet()

	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, z := range toCreate {
		if err = c.GetStore().CreateOrUpdate(
			z,
			object_mode.ModeApplyProto,
		); err != nil {
			// TODO-P2 add file for error handling
			c.handleStoreError(z, "", err)
			err = nil
			continue
		}

		results.Add(z)
	}

	if err = c.Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c CreateFromShas) handleStoreError(
	z *sku.Transacted,
	f string,
	in error,
) {
	var err error

	var normalError errors.StackTracer

	if errors.As(in, &normalError) {
		ui.Err().Printf("%s", normalError.Error())
	} else {
		err = errors.Errorf("writing zettel failed: %s: %s", f, in)
		ui.Err().Print(err)
	}
}
