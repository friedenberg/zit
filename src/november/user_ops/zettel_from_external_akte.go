package user_ops

import (
	"io"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/store_objekten"
	"github.com/friedenberg/zit/src/mike/umwelt"
)

type ZettelFromExternalAkte struct {
	*umwelt.Umwelt
	Etiketten etikett.Set
	Filter    script_value.ScriptValue
	Delete    bool
}

func (c ZettelFromExternalAkte) Run(
	ctx *errors.Ctx,
	args ...string,
) (results zettel_transacted.Set) {
	if ctx.Err = c.Lock(); !ctx.IsEmpty() {
		ctx.Wrap()
		return
	}

	defer ctx.Defer(c.Unlock)

	results = zettel_transacted.MakeSetUnique(len(args))

	for _, arg := range args {
		var z zettel.Zettel
		var tz zettel_transacted.Zettel

		if z, ctx.Err = c.zettelForAkte(arg); !ctx.IsEmpty() {
			ctx.Wrap()
			return
		}

		if tz, ctx.Err = c.StoreObjekten().Create(z); !ctx.IsEmpty() {
			ctx.Wrap()
			return
		}

		akteSha := tz.Named.Stored.Zettel.Akte

		if ctx.Err = c.StoreObjekten().AkteExists(akteSha); !ctx.IsEmpty() {
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

		results.Add(tz)

		if c.Delete {
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
	aktePath string,
) (z zettel.Zettel, err error) {
	var r io.Reader

	errors.Print("running")

	if r, err = c.Filter.Run(aktePath); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer c.Filter.Close()

	z.Etiketten = c.Etiketten

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

	z.Akte = akteWriter.Sha()

	if err = z.Bezeichnung.Set(path.Base(aktePath)); err != nil {
		err = errors.Wrap(err)
		return
	}

	ext := path.Ext(aktePath)

	if ext != "" {
		if err = z.Typ.Set(path.Ext(aktePath)); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
