package user_ops

import (
	"io"
	"log"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
	"github.com/friedenberg/zit/hotel/zettels"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type CreateFromPaths struct {
	// Options _ZettelsCheckinOptions
	Umwelt              *umwelt.Umwelt
	Format              _ZettelFormatsText
	Filter              _ScriptValue
	ReadHinweisFromPath bool
}

type CreateFromPathsResults struct {
	Zettelen []_Zettel
}

func (c CreateFromPaths) Run(args ...string) (results CreateFromPathsResults, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	toCreate := make([]stored_zettel.External, 0, len(args))

	for _, arg := range args {
		var toAdd []stored_zettel.External

		if toAdd, err = c.zettelsFromPath(store, arg); err != nil {
			err = _Errorf("zettel text format error for path: %s: %w", arg, err)
			return
		}

		toCreate = append(toCreate, toAdd...)
	}

	for _, z := range toCreate {
		var named _NamedZettel
		if c.ReadHinweisFromPath {
			head, tail := id.HeadTailFromFileName(z.Path)

			var h hinweis.Hinweis

			if h, err = hinweis.MakeBlindHinweis(head + "/" + tail); err != nil {
				err = _Error(err)
				return
			}

			if named, err = store.Zettels().CreateWithHinweis(z.Zettel, h); err != nil {
				//TODO add file for error handling
				c.handleStoreError(named, "", err)
				err = nil
				return
			}
		} else {
			if named, err = store.Zettels().Create(z.Zettel); err != nil {
				//TODO add file for error handling
				c.handleStoreError(named, "", err)
				err = nil
				return
			}
		}

		_Outf("[%s %s]\n", named.Hinweis, named.Sha)
	}

	return
}

func (c CreateFromPaths) zettelsFromPath(store store_with_lock.Store, p string) (out []stored_zettel.External, err error) {
	var r io.Reader

	log.Print("running")

	if r, err = c.Filter.Run(p); err != nil {
		err = _Error(err)
		return
	}

	defer c.Filter.Close()

	ctx := zettel.FormatContextRead{
		In:                r,
		AkteWriterFactory: store.Zettels(),
	}

	if _, err = c.Format.ReadFrom(&ctx); err != nil {
		err = _Error(err)
		return
	}

	if ctx.RecoverableError != nil {
		var errAkteInlineAndFilePath zettel_formats.ErrHasInlineAkteAndFilePath

		if errors.As(ctx.RecoverableError, &errAkteInlineAndFilePath) {
			var z1 _Zettel

			if z1, err = errAkteInlineAndFilePath.Recover(); err != nil {
				err = _Error(err)
				return
			}

			out = append(
				out,
				stored_zettel.External{
					Path:   p,
					Zettel: z1,
				},
			)
		} else {
			err = _Errorf("unsupported recoverable error: %w", ctx.RecoverableError)
			return
		}
	}

	out = append(
		out,
		stored_zettel.External{
			Path:   p,
			Zettel: ctx.Zettel,
		},
	)

	return
}

func (c CreateFromPaths) handleStoreError(z _NamedZettel, f string, in error) {
	var err error

	var lostError zettels.VerlorenAndGefundenError
	var normalError errors.StackTracer

	if errors.As(in, &lostError) {
		var p string

		if p, err = lostError.AddToLostAndFound(c.Umwelt.DirZit("Verloren+Gefunden")); err != nil {
			stdprinter.Error(err)
			return
		}

		_Outf("lost+found: %s: %s\n", lostError.Error(), p)

	} else if errors.As(in, &normalError) {
		stdprinter.Errf("%s\n", normalError.Error())
	} else {
		err = _Errorf("writing zettel failed: %s: %w", f, in)
		stdprinter.Error(err)
	}
}
