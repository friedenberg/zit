package user_ops

import (
	"os"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/files"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type ReadCheckedOut struct {
	Umwelt  *umwelt.Umwelt
	Options _ZettelsCheckinOptions
}

type ReadCheckedOutResults struct {
	Zettelen map[hinweis.Hinweis]stored_zettel.CheckedOut
}

func (op ReadCheckedOut) Run(paths ...string) (results ReadCheckedOutResults, err error) {
	results.Zettelen = make(map[hinweis.Hinweis]stored_zettel.CheckedOut)

	var store store_with_lock.Store

	if store, err = store_with_lock.New(op.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	for _, p := range paths {
		if op.Options.AddMdExtension {
			p = p + ".md"
		}

		checked_out := stored_zettel.CheckedOut{}

		checked_out.External, err = op.readExternal(store, p)

		if op.Options.IgnoreMissingHinweis && errors.Is(os.ErrNotExist, err) {
			err = nil
			//results.Zettelen[ez.Hinweis] = stored_zettel.External{}
			continue
		} else if err != nil {
			err = errors.Error(err)
			return
		}

		if checked_out.Internal, err = store.Zettels().Read(checked_out.External.Hinweis); err != nil {
			err = errors.Error(err)
			return
		}

		results.Zettelen[checked_out.External.Hinweis] = checked_out
	}

	return
}

func (op ReadCheckedOut) readExternal(store store_with_lock.Store, p string) (ez stored_zettel.External, err error) {
	ez.Path = p

	head, tail := id.HeadTailFromFileName(p)

	if ez.Hinweis, err = hinweis.MakeBlindHinweis(head + "/" + tail); err != nil {
		err = _Error(err)
		return
	}

	c := zettel.FormatContextRead{
		AkteWriterFactory: store.Zettels(),
	}

	var f *os.File

	if !files.Exists(p) {
		err = os.ErrNotExist
		return
	}

	if f, err = os.Open(p); err != nil {
		err = _Error(err)
		return
	}

	defer open_file_guard.Close(f)

	c.In = f

	if _, err = op.Options.Format.ReadFrom(&c); err != nil {
		err = _Errorf("%s: %w", f.Name(), err)
		return
	}

	ez.Zettel = c.Zettel
	ez.AktePath = c.AktePath

	return
}
