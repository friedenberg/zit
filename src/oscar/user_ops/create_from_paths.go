package user_ops

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/external"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/hotel/transacted"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/zettel_external"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type CreateFromPaths struct {
	*umwelt.Umwelt
	TextParser  metadatei.TextParser
	Filter      script_value.ScriptValue
	ProtoZettel zettel.ProtoZettel
	Delete      bool
	Dedupe      bool
	// ReadHinweisFromPath bool
}

func (c CreateFromPaths) Run(
	args ...string,
) (results schnittstellen.MutableSetLike[*transacted.Zettel], err error) {
	// TODO-P3 support different modes of de-duplication
	// TODO-P3 support merging of duplicated akten
	toCreate := zettel_external.MakeMutableSetUniqueFD()
	toDelete := zettel_external.MakeMutableSetUniqueFD()

	for _, arg := range args {
		if err = c.zettelsFromPath(
			arg,
			func(z *external.Zettel) (err error) {
				toCreate.Add(z)
				if c.Delete {
					toDelete.Add(z)
				}

				return
			},
		); err != nil {
			err = errors.Errorf(
				"zettel text format error for path: %s: %s",
				arg,
				err,
			)
			return
		}
	}

	results = zettel.MakeMutableSetHinweis(0)

	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, c.Unlock)

	if c.Dedupe {
		matcher := zettel_external.MakeMutableMatchSet(toCreate)

		if err = c.StoreObjekten().Zettel().ReadAll(
			iter.MakeChain(
				matcher.Match,
				iter.AddClone[transacted.Zettel, *transacted.Zettel](results),
			),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	err = results.Each(
		func(z *transacted.Zettel) (err error) {
			if c.ProtoZettel.Apply(z) {
				var zt *transacted.Zettel

				if zt, err = c.StoreObjekten().Zettel().Update(
					z,
					&z.Kennung,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				z = zt
			}

			return
		},
	)

	if err = toCreate.Each(
		func(z *external.Zettel) (err error) {
			if z.GetMetadatei().IsEmpty() {
				return
			}

			cz := zettel.CheckedOut{
				External: *z,
			}

			var zt *transacted.Zettel

			if zt, err = c.StoreObjekten().Zettel().Create(z); err != nil {
				// TODO-P2 add file for error handling
				c.handleStoreError(cz, "", err)
				err = nil
				return
			}

			cz.Internal = *zt

			if c.ProtoZettel.Apply(&cz.Internal) {
				if zt, err = c.StoreObjekten().Zettel().Update(
					cz.Internal,
					&cz.Internal.Kennung,
				); err != nil {
					// TODO-P2 add file for error handling
					c.handleStoreError(cz, "", err)
					err = nil
					return
				}

				cz.Internal = *zt
			}

			// TODO-P4 get matches
			cz.DetermineState(true)

			zv := &transacted.Zettel{}

			zv.ResetWith(cz.Internal)

			results.Add(zv)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = toDelete.Each(
		func(z *external.Zettel) (err error) {
			// TODO-P2 move to checkout store
			if err = os.Remove(z.GetObjekteFD().Path); err != nil {
				err = errors.Wrap(err)
				return
			}

			pathRel := c.Standort().RelToCwdOrSame(z.GetObjekteFD().Path)

			// TODO-P2 move to printer
			errors.Out().Printf("[%s] (deleted)", pathRel)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO-P1 migrate this to use store_working_directory
func (c *CreateFromPaths) zettelsFromPath(
	p string,
	wf schnittstellen.FuncIter[*external.Zettel],
) (err error) {
	var r io.Reader

	errors.Log().Print("running")

	if r, err = c.Filter.Run(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, &c.Filter)

	ze := &sku.External[kennung.Hinweis, *kennung.Hinweis]{
		FDs: sku.ExternalFDs{
			Objekte: kennung.FD{
				Path: p,
			},
		},
	}

	if _, err = c.TextParser.ParseMetadatei(r, ze); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.StoreObjekten().Zettel().SaveObjekte(ze); err != nil {
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
	z zettel.CheckedOut,
	f string,
	in error,
) {
	var err error

	var lostError objekte_store.VerlorenAndGefundenError
	var normalError errors.StackTracer

	if errors.As(in, &lostError) {
		var p string

		if p, err = lostError.AddToLostAndFound(c.Standort().DirZit("Verloren+Gefunden")); err != nil {
			errors.Err().Print(err)
			return
		}

		errors.Out().Printf("lost+found: %s: %s", lostError.Error(), p)

	} else if errors.As(in, &normalError) {
		errors.Err().Printf("%s", normalError.Error())
	} else {
		err = errors.Errorf("writing zettel failed: %s: %s", f, in)
		errors.Err().Print(err)
	}
}
