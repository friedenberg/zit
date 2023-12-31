package user_ops

import (
	"io"
	"os"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/india/objekte_collections"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type ZettelFromExternalAkte struct {
	*umwelt.Umwelt
	ProtoZettel zettel.ProtoZettel
	Filter      script_value.ScriptValue
	Delete      bool
	Dedupe      bool
}

func (c ZettelFromExternalAkte) Run(
	ms matcher.Query,
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

	for _, fd := range iter.SortedValues[*fd.FD](ms.GetCwdFDs()) {
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

		if err = c.StoreObjekten().ReadAll(
			gattungen.MakeSet(gattung.Zettel),
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
			if c.ProtoZettel.Apply(z) {
				if _, err = c.StoreObjekten().CreateOrUpdateTransacted(
					z,
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

		if tz, err = c.StoreObjekten().Create(z); err != nil {
			err = errors.Wrap(err)
			return
		}

		if c.ProtoZettel.Apply(tz) {
			if tz, err = c.StoreObjekten().CreateOrUpdateTransacted(
				tz,
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

	dp := c.Umwelt.PrinterFDDeleted()

	err = toDelete.Each(
		func(z *sku.External) (err error) {
			// TODO-P4 move to checkout store
			if err = os.Remove(z.GetAkteFD().GetPath()); err != nil {
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
