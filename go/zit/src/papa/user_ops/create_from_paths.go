package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/read_write_repo_local"
)

type CreateFromPaths struct {
	*read_write_repo_local.Repo
	sku.Proto
	TextParser object_metadata.TextParser
	Filter     script_value.ScriptValue
	Delete     bool
	// ReadHinweisFromPath bool
}

func (c CreateFromPaths) Run(
	args ...string,
) (results sku.TransactedMutableSet, err error) {
	toCreate := make(map[sha.Bytes]*sku.Transacted)
	toDelete := fd.MakeMutableSet()

	o := sku.CommitOptions{
		StoreOptions: sku.GetStoreOptionsRealizeWithProto(),
	}

	for _, arg := range args {
		var z *sku.Transacted
		var i sku.FSItem

		i.Reset()

		i.ExternalObjectId.SetGenre(genres.Zettel)

		if err = i.Object.Set(arg); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = i.Add(&i.Object); err != nil {
			err = errors.Wrap(err)
			return
		}

		if z, err = c.GetStore().GetStoreFS().ReadExternalFromItem(
			o,
			&i,
			nil,
		); err != nil {
			err = errors.Errorf(
				"zettel text format error for path: %s: %s",
				arg,
				err,
			)
			return
		}

		sh := &z.Metadata.Shas.SelfMetadataWithoutTai

		if sh.IsNull() {
			return
		}

		k := sh.GetBytes()
		existing, ok := toCreate[k]

		if ok {
			if err = existing.Metadata.Description.Set(
				z.Metadata.Description.String(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else {
			toCreate[k] = z
		}

		if c.Delete {
			{
				var object *fd.FD

				if object, err = c.GetStore().GetStoreFS().GetObjectOrError(z); err != nil {
					err = errors.Wrap(err)
					return
				}

				var f fd.FD
				f.ResetWith(object)
				toDelete.Add(&f)
			}

			{
				var blob *fd.FD

				if blob, err = c.GetStore().GetStoreFS().GetObjectOrError(z); err != nil {
					err = errors.Wrap(err)
					return
				}

				var f fd.FD
				f.ResetWith(blob)
				toDelete.Add(&f)
			}
		}
	}

	results = sku.MakeTransactedMutableSet()

	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, z := range toCreate {
		if z.Metadata.IsEmpty() {
			return
		}

		if err = c.GetStore().CreateOrUpdate(
			z,
			sku.StoreOptions{
				ApplyProto: true,
			},
		); err != nil {
			// TODO-P2 add file for error handling
			c.handleStoreError(z, "", err)
			err = nil
			continue
		}

		results.Add(z)
	}

	if err = toDelete.Each(
		func(f *fd.FD) (err error) {
			// TODO-P2 move to checkout store
			if err = c.GetRepoLayout().Delete(f.GetPath()); err != nil {
				err = errors.Wrap(err)
				return
			}

			pathRel := c.GetRepoLayout().RelToCwdOrSame(f.GetPath())

			// TODO-P2 move to printer
			c.GetUI().Printf("[%s] (deleted)", pathRel)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c CreateFromPaths) handleStoreError(
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
