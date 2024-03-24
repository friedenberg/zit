package user_ops

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/charlie/script_value"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/objekte_collections"
	"code.linenisgreat.com/zit/src/india/query"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/src/kilo/objekte_store"
	"code.linenisgreat.com/zit/src/kilo/zettel"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
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
) (results sku.TransactedMutableSet, err error) {
	// TODO-P3 support different modes of de-duplication
	// TODO-P3 support merging of duplicated akten
	toCreate := objekte_collections.MakeMutableSetUniqueFD()
	toDelete := objekte_collections.MakeMutableSetUniqueFD()

	for _, arg := range args {
		if err = c.zettelsFromPath(
			arg,
			func(z *sku.External) (err error) {
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

	results = sku.MakeTransactedMutableSet()

	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, c.Unlock)

	if c.Dedupe {
		matcher := objekte_collections.MakeMutableMatchSet(toCreate)

		b := c.MakeMetaIdSetWithoutExcludedHidden(kennung.MakeGattung(gattung.Zettel))

		var qg *query.Group

		if qg, err = b.BuildQueryGroup(); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = c.Store().QueryWithoutCwd(
			qg,
			iter.MakeChain(
				matcher.Match,
				func(sk *sku.Transacted) (err error) {
					var z sku.Transacted

					if err = z.SetFromSkuLike(sk); err != nil {
						err = errors.Wrap(err)
						return
					}

					return results.Add(&z)
				},
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
		func(z *sku.Transacted) (err error) {
			if c.ProtoZettel.Apply(z) {
				if _, err = c.Store().CreateOrUpdateTransacted(z); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
	)

	if err = toCreate.Each(
		func(z *sku.External) (err error) {
			if z.Metadatei.IsEmpty() {
				return
			}

			cz := sku.CheckedOut{}

			var zt *sku.Transacted

			if zt, err = c.Store().Create(z); err != nil {
				// TODO-P2 add file for error handling
				c.handleStoreError(cz, "", err)
				err = nil
				return
			}

			if err = cz.External.Transacted.SetFromSkuLike(zt); err != nil {
				err = errors.Wrapf(err, "Sku: %q", sku_fmt.String(&z.Transacted))
				return
			}

			if err = cz.Internal.SetFromSkuLike(zt); err != nil {
				err = errors.Wrap(err)
				return
			}

			if c.ProtoZettel.Apply(&cz.Internal) {
				if zt, err = c.Store().CreateOrUpdateTransacted(
					&cz.Internal,
				); err != nil {
					// TODO-P2 add file for error handling
					c.handleStoreError(cz, "", err)
					err = nil
					return
				}

				if err = cz.Internal.SetFromSkuLike(zt); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			// TODO-P4 get matches
			cz.DetermineState(true)

			zv := &sku.Transacted{}

			if err = zv.Kennung.SetWithKennung(&kennung.Hinweis{}); err != nil {
				err = errors.Wrap(err)
				return
			}

			sku.TransactedResetter.ResetWith(zv, &cz.Internal)

			results.Add(zv)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = toDelete.Each(
		func(z *sku.External) (err error) {
			// TODO-P2 move to checkout store
			if err = os.Remove(z.GetObjekteFD().GetPath()); err != nil {
				err = errors.Wrap(err)
				return
			}

			pathRel := c.Standort().RelToCwdOrSame(z.GetObjekteFD().GetPath())

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
	wf schnittstellen.FuncIter[*sku.External],
) (err error) {
	var r io.Reader

	errors.Log().Print("running")

	if r, err = c.Filter.Run(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, &c.Filter)

	var fd fd.FD

	if err = fd.SetPath(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	ze := sku.GetExternalPool().Get()
	ze.FDs = sku.ExternalFDs{
		Objekte: fd,
	}

	if err = ze.Kennung.SetWithKennung(&kennung.Hinweis{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = c.TextParser.ParseMetadatei(r, ze); err != nil {
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
	z sku.CheckedOut,
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
