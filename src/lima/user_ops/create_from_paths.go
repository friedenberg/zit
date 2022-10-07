package user_ops

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/echo/zettel_stored"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/golf/zettel_external"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/india/store_objekten"
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

type CreateFromPaths struct {
	*umwelt.Umwelt
	Format zettel.Format
	Filter script_value.ScriptValue
	Delete bool
	// ReadHinweisFromPath bool
}

func (c CreateFromPaths) Run(args ...string) (results zettel_checked_out.Set, err error) {
	toCreate := make([]zettel_external.Zettel, 0, len(args))

	for _, arg := range args {
		var toAdd []zettel_external.Zettel

		if toAdd, err = c.zettelsFromPath(arg); err != nil {
			err = errors.Errorf("zettel text format error for path: %s: %s", arg, err)
			return
		}

		toCreate = append(toCreate, toAdd...)
	}

	results = zettel_checked_out.MakeSetUnique(len(toCreate))

	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer c.Unlock()

	for _, z := range toCreate {
		cz := zettel_checked_out.Zettel{
			External: z,
		}

		//TODO
		if false /*c.ReadHinweisFromPath*/ {
			//head, tail := id.HeadTailFromFileName(z.ZettelFD.Path)

			//var h hinweis.Hinweis

			//if h, err = hinweis.Make(head + "/" + tail); err != nil {
			//	err = errors.Error(err)
			//	return
			//}

			//if tz, err = store.StoreObjekten().CreateWithHinweis(z.Stored.Zettel, h); err != nil {
			//	//TODO add file for error handling
			//	c.handleStoreError(tz, "", err)
			//	err = nil
			//	return
			//}
		} else {
			if cz.Internal, err = c.StoreObjekten().Create(z.Named.Stored.Zettel); err != nil {
				//TODO add file for error handling
				c.handleStoreError(cz, "", err)
				err = nil
				return
			}

			//TODO get matches
			cz.DetermineState()

			results.Add(cz)

			if c.Delete {
				//TODO move to checkout store
				if err = os.Remove(cz.External.ZettelFD.Path); err != nil {
					err = errors.Wrap(err)
					return
				}

				//TODO move to printer
				errors.PrintOutf("[%s] (deleted)", cz.External.ZettelFD.Path)
			}
		}
	}

	return
}

func (c CreateFromPaths) zettelsFromPath(p string) (out []zettel_external.Zettel, err error) {
	var r io.Reader

	errors.Print("running")

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

	if ctx.RecoverableError != nil {
		var errAkteInlineAndFilePath zettel.ErrHasInlineAkteAndFilePath

		if errors.As(ctx.RecoverableError, &errAkteInlineAndFilePath) {
			var z1 zettel.Zettel

			if z1, err = errAkteInlineAndFilePath.Recover(); err != nil {
				err = errors.Wrap(err)
				return
			}

			out = append(
				out,
				zettel_external.Zettel{
					ZettelFD: zettel_external.FD{
						Path: p,
					},
					Named: zettel_named.Zettel{
						Stored: zettel_stored.Stored{
							//TODO sha?
							Zettel: z1,
						},
					},
				},
			)
		} else {
			err = errors.Errorf("unsupported recoverable error: %s", ctx.RecoverableError)
			return
		}
	}

	var s sha.Sha

	if s, err = ctx.Zettel.ObjekteSha(); err != nil {
		err = errors.Wrap(err)
		return
	}

	out = append(
		out,
		zettel_external.Zettel{
			ZettelFD: zettel_external.FD{
				Path: p,
			},
			Named: zettel_named.Zettel{
				Stored: zettel_stored.Stored{
					Sha:    s,
					Zettel: ctx.Zettel,
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
			errors.PrintErr(err)
			return
		}

		errors.PrintOutf("lost+found: %s: %s", lostError.Error(), p)

	} else if errors.As(in, &normalError) {
		errors.PrintErrf("%s", normalError.Error())
	} else {
		err = errors.Errorf("writing zettel failed: %s: %s", f, in)
		errors.PrintErr(err)
	}
}
