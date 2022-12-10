package user_ops

import (
	"io"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/delta/sha"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/india/zettel"
	"github.com/friedenberg/zit/src/india/zettel_external"
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
	args ...string,
) (results zettel.MutableSet, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer c.Unlock()

	toCreate := zettel_external.MakeMutableSetUniqueAkte()
	toDelete := zettel_external.MakeMutableSetUniqueFD()

	results = zettel.MakeMutableSetUnique(len(args))

	for _, arg := range args {
		var z *zettel_external.Zettel

		akteFD := fd.FD{
			Path: arg,
		}

		if z, err = c.zettelForAkte(akteFD); err != nil {
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

		if err = c.StoreObjekten().Zettel().ReadAllTransacted(
			matcher.Match,
			results.AddAndDoNotRepool,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	err = results.Each(
		func(z *zettel.Transacted) (err error) {
			if c.ProtoZettel.Apply(&z.Objekte) {
				if *z, err = c.StoreObjekten().Zettel().Update(
					&z.Objekte,
					&z.Sku.Kennung,
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	err = toCreate.Each(
		func(z *zettel_external.Zettel) (err error) {
			if z.Objekte.IsEmpty() {
				return
			}

			var tz zettel.Transacted

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

			results.Add(&tz)

			return
		},
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	err = toDelete.Each(
		func(z *zettel_external.Zettel) (err error) {
			//TODO-P4 move to checkout store
			if err = os.Remove(z.AkteFD.Path); err != nil {
				err = errors.Wrap(err)
				return
			}

			pathRel := c.Standort().RelToCwdOrSame(z.AkteFD.Path)

			//TODO-P4 move to printer
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
	akteFD fd.FD,
) (z *zettel_external.Zettel, err error) {
	z = &zettel_external.Zettel{
		AkteFD: akteFD,
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

	z.Objekte.Reset(nil)
	z.Objekte.Akte = akteWriter.Sha()

	//TODO-P4 move to protozettel
	if err = z.Objekte.Bezeichnung.Set(
		path.Base(akteFD.Path),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO-P4 use konfig
	ext := akteFD.Ext()

	if ext != "" {
		if err = z.Objekte.Typ.Set(akteFD.Ext()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
