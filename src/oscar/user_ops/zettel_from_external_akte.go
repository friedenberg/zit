package user_ops

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/juliett/zettel"
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
	ms kennung.MetaSet,
) (results zettel.MutableSet, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, c.Unlock)

	toCreate := zettel_external.MakeMutableSetUniqueAkte()
	toDelete := zettel_external.MakeMutableSetUniqueFD()

	results = zettel.MakeMutableSetHinweis(0)

	fds := collections.SortedValues(ms.GetFDs())

	for _, fd := range fds {
		var z *zettel.External

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
				collections.AddClone[zettel.Transacted, *zettel.Transacted](results),
			),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = results.Each(
		func(z *zettel.Transacted) (err error) {
			if c.ProtoZettel.Apply(&z.Objekte) {
				if z, err = c.StoreObjekten().Zettel().Update(
					&z.Objekte,
					&z.Sku.Kennung,
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

	err = toCreate.Each(
		func(z *zettel.External) (err error) {
			if z.Objekte.IsEmpty() {
				return
			}

			var tz *zettel.Transacted

			if tz, err = c.StoreObjekten().Zettel().Create(z.Objekte); err != nil {
				err = errors.Wrap(err)
				return
			}

			if c.ProtoZettel.Apply(&tz.Objekte) {
				if tz, err = c.StoreObjekten().Zettel().Update(
					&tz.Objekte,
					&tz.Sku.Kennung,
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			results.Add(tz)

			return
		},
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	err = toDelete.Each(
		func(z *zettel.External) (err error) {
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

func (c ZettelFromExternalAkte) zettelForAkte(
	akteFD kennung.FD,
) (z *zettel.External, err error) {
	z = &zettel.External{
		Sku: sku.External[kennung.Hinweis, *kennung.Hinweis]{
			FDs: sku.ExternalFDs{
				Akte: akteFD,
			},
		},
	}

	var r io.Reader

	errors.Log().Print("running")

	if r, err = c.Filter.Run(akteFD.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer c.Filter.Close()

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

	z.Objekte.Reset()
	akteSha := sha.Make(akteWriter.Sha())
	z.Objekte.Akte = akteSha
	z.Sku.AkteSha = akteSha

	if err = c.ProtoZettel.ApplyWithAkteFD(&z.Objekte, akteFD); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
