package user_ops

import (
	"os"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/akte_ext"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

type WriteEmptyZettel struct {
	// Options _ZettelsCheckinOptions
	Umwelt _Umwelt
	Store  _Store
	Format zettel.Format
	Filter _ScriptValue
}

type WriteEmptyZettelResults struct {
	Zettel stored_zettel.External
}

func (c WriteEmptyZettel) Run() (results WriteEmptyZettelResults, err error) {
	var hinweis hinweis.Hinweis

	if hinweis, err = c.Store.Hinweisen().Factory().Make(); err != nil {
		err = errors.Error(err)
		return
	}

	var dir string

	if dir, err = os.Getwd(); err != nil {
		err = _Error(err)
		return
	}

	var filename string

	if filename, err = id.MakeDirIfNecessary(hinweis, dir); err != nil {
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

	results.Zettel.Path = f.Name()

	etiketten := etikett.NewSet()

	for e, t := range c.Umwelt.Konfig.Tags {
		if !t.AddToNewZettels {
			continue
		}

		if err = etiketten.AddString(e); err != nil {
			err = _Error(err)
			return
		}
	}

	results.Zettel.Zettel = zettel.Zettel{
		Etiketten: etiketten,
		AkteExt:   akte_ext.AkteExt{Value: "md"},
	}

	ctx := zettel.FormatContextWrite{
		Out:               f,
		AkteReaderFactory: c.Store,
		Zettel:            results.Zettel.Zettel,
	}

	if _, err = c.Format.WriteTo(ctx); err != nil {
		err = _Error(err)
		return
	}

	return
}
