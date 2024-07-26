package user_ops

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type CreateFromPaths struct {
	*env.Env
	sku.Proto
	TextParser object_metadata.TextParser
	Filter     script_value.ScriptValue
	Delete     bool
	// ReadHinweisFromPath bool
}

func (c CreateFromPaths) Run(
	args ...string,
) (results sku.TransactedMutableSet, err error) {
	toCreate := make(map[sha.Bytes]*store_fs.External)
	toDelete := fd.MakeMutableSet()

	o := sku.CommitOptions{
		Mode: objekte_mode.ModeRealizeWithProto,
	}

	for _, arg := range args {
		var z *store_fs.External
		var t store_fs.ObjectIdFDPair

		t.ObjectId.SetGenre(genres.Zettel)

		if err = t.FDs.Object.Set(arg); err != nil {
			err = errors.Wrap(err)
			return
		}

		if z, err = c.GetStore().GetCwdFiles().ReadExternalFromObjectIdFDPair(
			o,
			&t,
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
				var f fd.FD
				f.ResetWith(&z.FDs.Object)
				toDelete.Add(&f)
			}

			{
				var f fd.FD
				f.ResetWith(&z.FDs.Blob)
				toDelete.Add(&f)
			}
		}
	}

	results = sku.MakeTransactedMutableSet()

	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, c.Unlock)

	for _, z := range toCreate {
		if z.Metadata.IsEmpty() {
			return
		}

		if err = c.GetStore().CreateOrUpdate(
			&z.Transacted,
			objekte_mode.ModeApplyProto,
		); err != nil {
			// TODO-P2 add file for error handling
			c.handleStoreError(z, "", err)
			err = nil
			continue
		}

		results.Add(&z.Transacted)
	}

	if err = toDelete.Each(
		func(f *fd.FD) (err error) {
			// TODO-P2 move to checkout store
			if err = c.GetFSHome().Delete(f.GetPath()); err != nil {
				err = errors.Wrap(err)
				return
			}

			pathRel := c.GetFSHome().RelToCwdOrSame(f.GetPath())

			// TODO-P2 move to printer
			ui.Out().Printf("[%s] (deleted)", pathRel)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO-P1 migrate this to use store_working_directory
// TODO remove this
func (c *CreateFromPaths) zettelsFromPath(
	p string,
	wf interfaces.FuncIter[*store_fs.External],
) (err error) {
	var r io.Reader

	ui.Log().Print("running")

	if r, err = c.Filter.Run(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, &c.Filter)

	var fd fd.FD

	if err = fd.Set(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	ze := store_fs.GetExternalPool().Get()
	ze.FDs = store_fs.FDPair{
		Object: fd,
	}

	ze.Metadata.Tai = ids.TaiFromTime(fd.ModTime())

	ze.ObjectId.SetGenre(genres.Zettel)

	if _, err = c.TextParser.ParseMetadata(r, ze); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = ze.CalculateObjectShas(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = wf(ze); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c CreateFromPaths) handleStoreError(
	z *store_fs.External,
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
