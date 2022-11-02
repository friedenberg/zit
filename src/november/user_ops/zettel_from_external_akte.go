package user_ops

import (
	"io"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_external"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
	"github.com/friedenberg/zit/src/mike/umwelt"
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
) (results zettel_transacted.MutableSet, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer c.Unlock()

	toCreate := zettel_external.MakeMutableSetUniqueAkte()
	toDelete := zettel_external.MakeMutableSetUniqueFD()

	results = zettel_transacted.MakeMutableSetUnique(len(args))

	for _, arg := range args {
		var z *zettel_external.Zettel

		akteFD := zettel_external.FD{
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

		writerMatches := zettel_transacted.MakeWriterZettelNamed(
			matcher.WriterZettelNamed(),
		)

		if err = c.StoreObjekten().ReadAllTransacted(
			writerMatches,
			zettel_transacted.MakeWriter(results.Add),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	err = results.Each(
		func(z *zettel_transacted.Zettel) (err error) {
			if c.ProtoZettel.Apply(&z.Named.Stored.Zettel) {
				if *z, err = c.StoreObjekten().Update(
					&z.Named,
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
			if z.Named.Stored.Zettel.IsEmpty() {
				return
			}

			var tz zettel_transacted.Zettel

			if tz, err = c.StoreObjekten().Create(z.Named.Stored.Zettel); err != nil {
				err = errors.Wrap(err)
				return
			}

			if c.ProtoZettel.Apply(&tz.Named.Stored.Zettel) {
				if tz, err = c.StoreObjekten().Update(&tz.Named); err != nil {
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
			//TODO move to checkout store
			if err = os.Remove(z.AkteFD.Path); err != nil {
				err = errors.Wrap(err)
				return
			}

			pathRel := c.Standort().RelToCwdOrSame(z.AkteFD.Path)

			//TODO move to printer
			errors.PrintOutf("[%s] (deleted)", pathRel)

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
	akteFD zettel_external.FD,
) (z *zettel_external.Zettel, err error) {
	z = &zettel_external.Zettel{
		AkteFD: akteFD,
	}

	var r io.Reader

	errors.Print("running")

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

	z.Named.Stored.Zettel.Akte = akteWriter.Sha()

	//TODO move to protozettel
	if err = z.Named.Stored.Zettel.Bezeichnung.Set(
		path.Base(akteFD.Path),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO use konfig
	ext := akteFD.Ext()

	if ext != "" {
		if err = z.Named.Stored.Zettel.Typ.Set(akteFD.Ext()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
