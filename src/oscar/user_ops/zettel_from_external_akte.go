package user_ops

import (
	"io"
	"os"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/external"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/india/transacted"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/kilo/zettel_external"
	"github.com/friedenberg/zit/src/november/umwelt"
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
) (results zettel.MutableSet, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, c.Unlock)

	toCreate := zettel_external.MakeMutableSetUniqueAkte()
	toDelete := zettel_external.MakeMutableSetUniqueFD()

	results = zettel.MakeMutableSetHinweis(0)

	fds := iter.SortedValues(ms.GetCwdFDs())

	for _, fd := range fds {
		var z *external.Zettel

		if z, err = c.zettelForAkte(fd); err != nil {
			err = errors.Wrap(err)
			return
		}

		toCreate.Add(z)

		if c.Delete {
			toDelete.Add(z)
		}
	}

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

	if err = results.Each(
		func(z *transacted.Zettel) (err error) {
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

		var tz *transacted.Zettel

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

		results.Add(tz)
	}

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	err = toDelete.Each(
		func(z *external.Zettel) (err error) {
			// TODO-P4 move to checkout store
			if err = os.Remove(z.GetAkteFD().Path); err != nil {
				err = errors.Wrap(err)
				return
			}

			pathRel := c.Standort().RelToCwdOrSame(z.GetAkteFD().Path)

			// TODO-P4 move to printer
			errors.Out().Printf("[%s] (deleted)", pathRel)

			return
		},
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *ZettelFromExternalAkte) zettelForAkte(
	akteFD kennung.FD,
) (z *external.Zettel, err error) {
	z = &external.Zettel{
		FDs: sku.ExternalFDs{
			Akte: akteFD,
		},
	}

	var r io.Reader

	errors.Log().Print("running")

	if r, err = c.Filter.Run(akteFD.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, &c.Filter)

	var akteWriter sha.WriteCloser

	if akteWriter, err = c.StoreObjekten().AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = io.Copy(akteWriter, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = akteWriter.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	z.GetMetadateiPtr().Reset()

	if err = c.ProtoZettel.ApplyWithAkteFD(z, akteFD); err != nil {
		err = errors.Wrap(err)
		return
	}

	z.SetAkteSha(akteWriter.GetShaLike())

	return
}
