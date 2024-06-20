package user_ops

import (
	"io"
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/objekte_collections"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/zettel"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type ZettelFromExternalAkte struct {
	*umwelt.Umwelt
	ProtoZettel zettel.ProtoZettel
	Filter      script_value.ScriptValue
	Delete      bool
	Dedupe      bool
}

func (c ZettelFromExternalAkte) Run(
	qg *query.Group,
) (results sku.TransactedMutableSet, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, c.Unlock)

	toCreate := objekte_collections.MakeMutableSetUniqueAkte()
	toDelete := objekte_collections.MakeMutableSetUniqueFD()

	results = sku.MakeTransactedMutableSet()

	fds := fd.MakeMutableSetSha()

	for _, fd := range iter.SortedValues[*fd.FD](qg.GetCwdFDs()) {
		if err = c.processOneFD(fd, fds.Add); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = fds.Each(
		func(fd *fd.FD) (err error) {
			var z *sku.External

			if z, err = c.zettelForAkte(fd); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = toCreate.Add(z); err != nil {
				err = errors.Wrap(err)
				return
			}

			if !c.Delete {
				return
			}

			if err = toDelete.Add(z); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.Dedupe {
		matcher := objekte_collections.MakeMutableMatchSet(toCreate)

		if err = c.GetStore().QueryWithoutCwd(
			qg,
			iter.MakeChain(
				matcher.Match,
				func(sk *sku.Transacted) (err error) {
					z := &sku.Transacted{}

					if err = z.SetFromSkuLike(sk); err != nil {
						err = errors.Wrap(err)
						return
					}

					return results.Add(z)
				},
			),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = results.Each(
		func(z *sku.Transacted) (err error) {
			if c.ProtoZettel.Apply(z, gattung.Zettel) {
				if err = c.GetStore().CreateOrUpdateTransacted(
					z,
					false,
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sortedToCreated := iter.Elements[*sku.External](toCreate)

	sort.Slice(
		sortedToCreated,
		func(i, j int) bool {
			return sortedToCreated[i].GetAkteFD().String() < sortedToCreated[j].GetAkteFD().String()
		},
	)

	for _, z := range sortedToCreated {
		if z.Metadatei.IsEmpty() {
			return
		}

		var tz *sku.Transacted

		if tz, err = c.GetStore().Create(z); err != nil {
			err = errors.Wrap(err)
			return
		}

		if c.ProtoZettel.Apply(tz, gattung.Zettel) {
			if err = c.GetStore().CreateOrUpdateTransacted(
				tz,
				false,
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		results.Add(tz)
	}

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO move to umwelt
	dp := c.Umwelt.PrinterFDDeleted()

	err = toDelete.Each(
		func(z *sku.External) (err error) {
			if err = c.Umwelt.Standort().Delete(z.GetAkteFD().GetPath()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return dp(&z.GetFDsPtr().Akte)
		},
	)
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *ZettelFromExternalAkte) processOneFD(
	fd *fd.FD,
	add schnittstellen.FuncIter[*fd.FD],
) (err error) {
	var r io.Reader

	if r, err = c.Filter.Run(fd.GetPath()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, &c.Filter)

	var akteWriter sha.WriteCloser

	if akteWriter, err = c.Standort().AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, akteWriter)

	if _, err = io.Copy(akteWriter, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	fd.SetShaLike(akteWriter.GetShaLike())

	if err = add(fd); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *ZettelFromExternalAkte) zettelForAkte(
	akteFD *fd.FD,
) (z *sku.External, err error) {
	z = sku.GetExternalPool().Get()

	z.FDs.Akte.ResetWith(akteFD)

	if err = z.Transacted.Kennung.SetWithKennung(&kennung.Hinweis{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.ProtoZettel.ApplyWithAkteFD(z, akteFD); err != nil {
		err = errors.Wrap(err)
		return
	}

	z.SetAkteSha(akteFD.GetShaLike())

	return
}
