package user_ops

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
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
) (results schnittstellen.MutableSet[*zettel.Transacted], err error) {
	// TODO support different modes of de-duplication
	// TODO support merging of duplicated akten
	toCreate := zettel_external.MakeMutableSetUniqueFD()
	toDelete := zettel_external.MakeMutableSetUniqueFD()

	for _, arg := range args {
		if err = c.zettelsFromPath(
			arg,
			func(z *zettel.External) (err error) {
				toCreate.Add(z)
				if c.Delete {
					toDelete.Add(z)
				}

				return
			},
		); err != nil {
			err = errors.Errorf("zettel text format error for path: %s: %s", arg, err)
			return
		}
	}

	results = collections.MakeMutableSet[*zettel.Transacted](
		func(zv *zettel.Transacted) string {
			if zv == nil {
				return ""
			}

			return zv.Sku.Kennung.String()
		},
	)

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
				collections.AddClone[zettel.Transacted, *zettel.Transacted](results),
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
		func(z *zettel.Transacted) (err error) {
			if c.ProtoZettel.Apply(z) {
				var zt *zettel.Transacted

				if zt, err = c.StoreObjekten().Zettel().Update(
					&z.Objekte,
					z,
					&z.Sku.Kennung,
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
		func(z *zettel.External) (err error) {
			if z.GetMetadatei().IsEmpty() {
				return
			}

			cz := zettel.CheckedOut{
				External: *z,
			}

			var zt *zettel.Transacted

			if zt, err = c.StoreObjekten().Zettel().Create(z.Objekte, z); err != nil {
				// TODO add file for error handling
				c.handleStoreError(cz, "", err)
				err = nil
				return
			}

			cz.Internal = *zt

			if c.ProtoZettel.Apply(&cz.Internal) {
				if zt, err = c.StoreObjekten().Zettel().Update(
					&cz.Internal.Objekte,
					cz.Internal,
					&cz.Internal.Sku.Kennung,
				); err != nil {
					// TODO add file for error handling
					c.handleStoreError(cz, "", err)
					err = nil
					return
				}

				cz.Internal = *zt
			}

			// TODO get matches
			cz.DetermineState()

			zv := &zettel.Transacted{}

			zv.ResetWith(cz.Internal)

			results.Add(zv)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = toDelete.Each(
		func(z *zettel.External) (err error) {
			// TODO move to checkout store
			if err = os.Remove(z.GetObjekteFD().Path); err != nil {
				err = errors.Wrap(err)
				return
			}

			pathRel := c.Standort().RelToCwdOrSame(z.GetObjekteFD().Path)

			// TODO move to printer
			errors.Out().Printf("[%s] (deleted)", pathRel)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO migrate this to use store_working_directory
func (c *CreateFromPaths) zettelsFromPath(
	p string,
	wf schnittstellen.FuncIter[*zettel.External],
) (err error) {
	var r io.Reader

	errors.Log().Print("running")

	if r, err = c.Filter.Run(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, &c.Filter)

	ze := &zettel.External{
		Sku: sku.External[kennung.Hinweis, *kennung.Hinweis]{
			FDs: sku.ExternalFDs{
				Objekte: kennung.FD{
					Path: p,
				},
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

func (c CreateFromPaths) handleStoreError(z zettel.CheckedOut, f string, in error) {
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
