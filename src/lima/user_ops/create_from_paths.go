package user_ops

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/delta/script_value"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_formats"
	"github.com/friedenberg/zit/src/golf/zettel_stored"
	"github.com/friedenberg/zit/src/india/store_objekten"
	"github.com/friedenberg/zit/src/india/zettel_checked_out"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
)

type CreateFromPaths struct {
	Umwelt *umwelt.Umwelt
	Format zettel.Format
	Filter script_value.ScriptValue
	Delete bool
	// ReadHinweisFromPath bool
}

func (c CreateFromPaths) Run(args ...string) (results zettel_checked_out.Set, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	toCreate := make([]zettel_stored.External, 0, len(args))

	for _, arg := range args {
		var toAdd []zettel_stored.External

		if toAdd, err = c.zettelsFromPath(store, arg); err != nil {
			err = errors.Errorf("zettel text format error for path: %s: %s", arg, err)
			return
		}

		toCreate = append(toCreate, toAdd...)
	}

	results = zettel_checked_out.MakeSetUnique(len(toCreate))

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
			if cz.Internal, err = store.StoreObjekten().Create(z.Stored.Zettel); err != nil {
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
					err = errors.Error(err)
					return
				}

				stdprinter.Outf("[%s] (deleted)\n", cz.External.ZettelFD.Path)
			}
		}

		stdprinter.Outf("%s (created)\n", cz.Internal.Named)
	}

	return
}

func (c CreateFromPaths) zettelsFromPath(store store_with_lock.Store, p string) (out []zettel_stored.External, err error) {
	var r io.Reader

	logz.Print("running")

	if r, err = c.Filter.Run(p); err != nil {
		err = errors.Error(err)
		return
	}

	defer c.Filter.Close()

	ctx := zettel.FormatContextRead{
		In:                r,
		AkteWriterFactory: store.StoreObjekten(),
	}

	if _, err = c.Format.ReadFrom(&ctx); err != nil {
		err = errors.Error(err)
		return
	}

	if ctx.RecoverableError != nil {
		var errAkteInlineAndFilePath zettel_formats.ErrHasInlineAkteAndFilePath

		if errors.As(ctx.RecoverableError, &errAkteInlineAndFilePath) {
			var z1 zettel.Zettel

			if z1, err = errAkteInlineAndFilePath.Recover(); err != nil {
				err = errors.Error(err)
				return
			}

			out = append(
				out,
				zettel_stored.External{
					ZettelFD: zettel_stored.FD{
						Path: p,
					},
					Named: zettel_stored.Named{
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

	out = append(
		out,
		zettel_stored.External{
			ZettelFD: zettel_stored.FD{
				Path: p,
			},
			Named: zettel_stored.Named{
				Stored: zettel_stored.Stored{
					//TODO sha?
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

		if p, err = lostError.AddToLostAndFound(c.Umwelt.DirZit("Verloren+Gefunden")); err != nil {
			stdprinter.Error(err)
			return
		}

		stdprinter.Outf("lost+found: %s: %s\n", lostError.Error(), p)

	} else if errors.As(in, &normalError) {
		stdprinter.Errf("%s\n", normalError.Error())
	} else {
		err = errors.Errorf("writing zettel failed: %s: %s", f, in)
		stdprinter.Error(err)
	}
}
