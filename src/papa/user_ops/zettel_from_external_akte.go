package user_ops

import (
	"io"
	"os"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/charlie/sha"
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

	fds := collections_value.MakeMutableSet[fd.FD](
		fd.KeyerSha{},
	)

	for _, fd := range iter.SortedValues[fd.FD](ms.GetCwdFDs()) {
		if err = c.processOneFD(fd, fds.Add); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = fds.Each(
		func(fd fd.FD) (err error) {
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

		if err = c.StoreObjekten().Zettel().ReadAll(
			iter.MakeChain(
				matcher.Match,
				func(sk *sku.Transacted) (err error) {
					z := &sku.Transacted{}

					if err = z.SetFromSkuLike(sk); err != nil {
						err = errors.Wrap(err)
						return
					}

					return results.AddPtr(z)
				},
			),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = results.EachPtr(
		func(z *sku.Transacted) (err error) {
			if c.ProtoZettel.Apply(z) {
				if z, err = c.StoreObjekten().Zettel().Update(
					z,
					&z.Kennung,
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

	sortedToCreated := toCreate.Elements()

	sort.Slice(
		sortedToCreated,
		func(i, j int) bool {
			return sortedToCreated[i].GetAkteFD().String() < sortedToCreated[j].GetAkteFD().String()
		},
	)

	for _, z := range sortedToCreated {
		if z.GetMetadatei().IsEmpty() {
			return
		}

		var tz *sku.Transacted

		if tz, err = c.StoreObjekten().Zettel().Create(z); err != nil {
			err = errors.Wrap(err)
			return
		}

		if c.ProtoZettel.Apply(tz) {
			if tz, err = c.StoreObjekten().Zettel().Update(
				tz,
				&tz.Kennung,
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		results.AddPtr(tz)
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
	fd fd.FD,
	add schnittstellen.FuncIter[fd.FD],
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
	akteFD fd.FD,
) (z *sku.External, err error) {
	z = &sku.External{
		FDs: sku.ExternalFDs{
			Akte: akteFD,
		},
	}

	z.Transacted.Kennung.KennungPtr = &kennung.Hinweis{}

	z.GetMetadateiPtr().Reset()

	if err = c.ProtoZettel.ApplyWithAkteFD(z, akteFD); err != nil {
		err = errors.Wrap(err)
		return
	}

	z.SetAkteSha(akteFD.GetShaLike())

	return
}
