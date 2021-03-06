package user_ops

import (
	"os"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type WriteNewZettels struct {
	Umwelt *umwelt.Umwelt
	Format zettel.Format
	Filter _ScriptValue
}

func (c WriteNewZettels) Run(zettelen ...zettel.Zettel) (results stored_zettel.SetExternal, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	var dir string

	if dir, err = os.Getwd(); err != nil {
		err = _Error(err)
		return
	}

	results = stored_zettel.MakeSetExternal()

	for _, z := range zettelen {
		var external stored_zettel.External

		if external, err = c.runOne(store, dir, z); err != nil {
			err = errors.Error(err)
			return
		}

		results[external.Hinweis.String()] = external
	}

	return
}

func (c WriteNewZettels) runOne(store store_with_lock.Store, dir string, z zettel.Zettel) (result stored_zettel.External, err error) {
	if result.Hinweis, err = store.Hinweisen().Factory().Make(); err != nil {
		err = errors.Error(err)
		return
	}

	var filename string

	if filename, err = id.MakeDirIfNecessary(result.Hinweis, dir); err != nil {
		err = _Error(err)
		return
	}

	filename = filename + ".md"

	var f *os.File

	if f, err = open_file_guard.Create(filename); err != nil {
		err = _Error(err)
		return
	}

	defer open_file_guard.Close(f)

	result.Path = f.Name()
	result.Zettel = z

	ctx := zettel.FormatContextWrite{
		Out:               f,
		AkteReaderFactory: store.Zettels(),
		Zettel:            result.Zettel,
	}

	if _, err = c.Format.WriteTo(ctx); err != nil {
		err = _Error(err)
		return
	}

	return
}
