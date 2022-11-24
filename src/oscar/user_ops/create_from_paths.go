package user_ops

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/lima/store_objekten"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type CreateFromPaths struct {
	*umwelt.Umwelt
	Format      zettel.Format
	Filter      script_value.ScriptValue
	ProtoZettel zettel.ProtoZettel
	Delete      bool
	Dedupe      bool
	// ReadHinweisFromPath bool
}

func (c CreateFromPaths) Run(args ...string) (results zettel_transacted.MutableSet, err error) {
	//TODO support different modes of de-duplication
	//TODO support merging of duplicated akten
	toCreate := zettel_external.MakeMutableSetUniqueFD()
	toDelete := zettel_external.MakeMutableSetUniqueFD()

	for _, arg := range args {
		if err = c.zettelsFromPath(
			arg,
			func(z *zettel_external.Zettel) (err error) {
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

	results = zettel_transacted.MakeMutableSetHinweis(toCreate.Len())

	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer c.Unlock()

	if c.Dedupe {
		matcher := zettel_external.MakeMutableMatchSet(toCreate)

		if err = c.StoreObjekten().ReadAllTransacted(
			zettel_transacted.MakeWriterZettelNamed(matcher.Match),
			results.AddAndDoNotRepool,
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
		func(z *zettel_transacted.Zettel) (err error) {
			if c.ProtoZettel.Apply(&z.Named.Stored.Objekte) {
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

	err = toCreate.Each(
		func(z *zettel_external.Zettel) (err error) {
			cz := zettel_checked_out.Zettel{
				External: *z,
			}

			if z.Named.Stored.Objekte.IsEmpty() {
				return
			}

			if cz.Internal, err = c.StoreObjekten().Create(z.Named.Stored.Objekte); err != nil {
				//TODO add file for error handling
				c.handleStoreError(cz, "", err)
				err = nil
				return
			}

			if c.ProtoZettel.Apply(&cz.Internal.Named.Stored.Objekte) {
				if cz.Internal, err = c.StoreObjekten().Update(
					&cz.Internal.Named,
				); err != nil {
					//TODO add file for error handling
					c.handleStoreError(cz, "", err)
					err = nil
					return
				}
			}

			//TODO get matches
			cz.DetermineState()

			results.Add(&cz.Internal)

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
			if err = os.Remove(z.ZettelFD.Path); err != nil {
				err = errors.Wrap(err)
				return
			}

			pathRel := c.Standort().RelToCwdOrSame(z.ZettelFD.Path)

			//TODO move to printer
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

// TODO migrate this to use store_working_directory
func (c CreateFromPaths) zettelsFromPath(
	p string,
	wf collections.WriterFunc[*zettel_external.Zettel],
) (err error) {
	var r io.Reader

	errors.Log().Print("running")

	if r, err = c.Filter.Run(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer c.Filter.Close()

	ctx := zettel.FormatContextRead{
		In:                r,
		AkteWriterFactory: c.StoreObjekten(),
	}

	if _, err = c.Format.ReadFrom(&ctx); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, e := range ctx.RecoverableErrors {
		var errAkteInlineAndFilePath zettel.ErrHasInlineAkteAndFilePath

		if errors.As(e, &errAkteInlineAndFilePath) {
			var z1 zettel.Zettel

			if z1, err = errAkteInlineAndFilePath.Recover(); err != nil {
				err = errors.Wrap(err)
				return
			}

			var s sha.Sha

			if s, err = z1.ObjekteSha(); err != nil {
				err = errors.Wrap(err)
				return
			}

			wf(
				&zettel_external.Zettel{
					ZettelFD: zettel_external.FD{
						Path: p,
					},
					Named: zettel.Named{
						Stored: zettel.Stored{
							Sha:     s,
							Objekte: z1,
						},
					},
				},
			)
		} else {
			err = errors.Errorf("unsupported recoverable error: %s", e)
			return
		}
	}

	var s sha.Sha

	if s, err = ctx.Zettel.ObjekteSha(); err != nil {
		err = errors.Wrap(err)
		return
	}

	wf(
		&zettel_external.Zettel{
			ZettelFD: zettel_external.FD{
				Path: p,
			},
			Named: zettel.Named{
				Stored: zettel.Stored{
					Sha:     s,
					Objekte: ctx.Zettel,
				},
			},
		},
	)

	return
}

func (c CreateFromPaths) handleStoreError(z zettel_checked_out.Zettel, f string, in error) {
	var err error

	var lostError store_objekten.VerlorenAndGefundenError
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
