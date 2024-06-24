package user_ops

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/objekte_collections"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/objekte"
	"code.linenisgreat.com/zit/go/zit/src/kilo/zettel"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type CreateFromPaths struct {
	*umwelt.Umwelt
	TextParser  metadatei.TextParser
	Filter      script_value.ScriptValue
	ProtoZettel zettel.ProtoZettel
	Delete      bool
	// ReadHinweisFromPath bool
}

func (c CreateFromPaths) Run(
	args ...string,
) (results sku.TransactedMutableSet, err error) {
	toCreate := make(map[sha.Bytes]*store_fs.External)
	toDelete := objekte_collections.MakeMutableSetUniqueFD()

	o := store.ObjekteOptions{
		Mode: objekte_mode.ModeRealizeWithProto,
	}

	for _, arg := range args {
		var z *store_fs.External
		var t store_fs.KennungFDPair

		t.Kennung.SetGattung(gattung.Zettel)

		if err = t.FDs.Objekte.Set(arg); err != nil {
			err = errors.Wrap(err)
			return
		}

		if z, err = c.GetStore().ReadOneExternalFS(
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

		sh := &z.Metadatei.Shas.SelbstMetadateiSansTai

		if sh.IsNull() {
			return
		}

		k := sh.GetBytes()
		existing, ok := toCreate[k]

		if ok {
			if err = existing.Metadatei.Bezeichnung.Set(
				z.Metadatei.Bezeichnung.String(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else {
			toCreate[k] = z
		}

		if c.Delete {
			toDelete.Add(z)
		}
	}

	results = sku.MakeTransactedMutableSet()

	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, c.Unlock)

	for _, z := range toCreate {
		if z.Metadatei.IsEmpty() {
			return
		}

		if err = c.GetStore().CreateOrUpdateTransacted(
			&z.Transacted,
			false,
		); err != nil {
			// TODO-P2 add file for error handling
			c.handleStoreError(z, "", err)
			err = nil
			continue
		}

		results.Add(&z.Transacted)
	}

	if err = toDelete.Each(
		func(z *store_fs.External) (err error) {
			// TODO-P2 move to checkout store
			if err = c.Standort().Delete(z.GetObjekteFD().GetPath()); err != nil {
				err = errors.Wrap(err)
				return
			}

			pathRel := c.Standort().RelToCwdOrSame(z.GetObjekteFD().GetPath())

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
	wf schnittstellen.FuncIter[*store_fs.External],
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
		Objekte: fd,
	}

	ze.Metadatei.Tai = kennung.TaiFromTime(fd.ModTime())

	ze.Kennung.SetGattung(gattung.Zettel)

	if _, err = c.TextParser.ParseMetadatei(r, ze); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = ze.CalculateObjekteShas(); err != nil {
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

	var lostError objekte.VerlorenAndGefundenError
	var normalError errors.StackTracer

	if errors.As(in, &lostError) {
		var p string

		if p, err = lostError.AddToLostAndFound(c.Standort().DirZit("Verloren+Gefunden")); err != nil {
			ui.Err().Print(err)
			return
		}

		ui.Out().Printf("lost+found: %s: %s", lostError.Error(), p)

	} else if errors.As(in, &normalError) {
		ui.Err().Printf("%s", normalError.Error())
	} else {
		err = errors.Errorf("writing zettel failed: %s: %s", f, in)
		ui.Err().Print(err)
	}
}
