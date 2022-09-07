package user_ops

import (
	"io"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/delta/umwelt"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/store_objekten"
	"github.com/friedenberg/zit/src/juliett/store_with_lock"
)

type ZettelFromExternalAkte struct {
	Umwelt    *umwelt.Umwelt
	Etiketten etikett.Set
	Delete    bool
}

func (c ZettelFromExternalAkte) Run(
	ctx *errors.Ctx,
	args ...string,
) (results zettel_transacted.Set) {
	var store store_with_lock.Store

	if store, ctx.Err = store_with_lock.New(c.Umwelt); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	defer ctx.Defer(store.Flush)

	results = zettel_transacted.MakeSetUnique(len(args))

	for _, arg := range args {
		var z zettel.Zettel
		var tz zettel_transacted.Zettel

		errors.PrintErr("creating zettel for akte")
		if z = c.zettelForAkte(ctx, store, arg); !ctx.IsEmpty() {
			errors.PrintErr("ctx error", ctx.IsEmpty())
			errors.PrintErr(ctx.Err)
			ctx.Wrap()
			return
		}

		errors.PrintErr("creating zettel")
		if tz, ctx.Err = store.StoreObjekten().Create(z); !ctx.IsEmpty() {
			ctx.Wrap()
			return
		}

		akteSha := tz.Named.Stored.Zettel.Akte

		if ctx.Err = store.StoreObjekten().AkteExists(akteSha); !ctx.IsEmpty() {
			if errors.Is(ctx.Err, store_objekten.ErrAkteExists{}) {
				err1 := ctx.Err.(store_objekten.ErrAkteExists)
				errors.PrintOutf("[%s %s] (has Akte matches)", arg, akteSha)
				err1.Set.Each(
					func(tz1 zettel_transacted.Zettel) (err error) {
						if tz1.Named.Hinweis.Equals(tz.Named.Hinweis) {
							return
						}
						//TODO eliminate zettels marked as duplicates / hidden
						errors.PrintOutf("\t%s", tz1.Named)
						return
					},
				)
				ctx.Err = nil
			} else {
				ctx.Wrapf("%s", arg)
				return
			}
		}

		errors.PrintErr("adding result")
		results.Add(tz)

		if c.Delete {
			errors.PrintErr("deleting old akte file")
			if ctx.Err = os.Remove(arg); !ctx.IsEmpty() {
				ctx.Wrap()
				return
			}

			errors.PrintErrf("[%s] (deleted)", arg)
		}

		//TODO-P3,D3 only emit if created rather than refound
		errors.PrintOutf("%s (created)", tz.Named)
	}

	return
}

func (c ZettelFromExternalAkte) zettelForAkte(
	ctx *errors.Ctx,
	store store_with_lock.Store,
	aktePath string,
) (z zettel.Zettel) {
	z.Etiketten = c.Etiketten

	var akteWriter sha.WriteCloser

	if akteWriter, ctx.Err = store.StoreObjekten().AkteWriter(); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	var f *os.File

	if f, ctx.Err = files.Open(aktePath); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	defer ctx.Defer(func() error { return files.Close(f) })

	if _, ctx.Err = io.Copy(akteWriter, f); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	if ctx.Err = akteWriter.Close(); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	z.Akte = akteWriter.Sha()

	if ctx.Err = z.Bezeichnung.Set(path.Base(aktePath)); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	if ctx.Err = z.Typ.Set(path.Ext(aktePath)); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	return
}
