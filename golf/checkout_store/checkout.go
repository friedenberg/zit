package checkout_store

import (
	"os"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/stdprinter"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

func (s *Store) Checkout(
	options CheckinOptions,
	zs ...stored_zettel.Transacted,
) (czs []stored_zettel.CheckedOut, err error) {
	czs = make([]stored_zettel.CheckedOut, len(zs))

	var dir string

	if dir, err = os.Getwd(); err != nil {
		err = errors.Error(err)
		return
	}

	for i, sz := range zs {
		var filename string

		if filename, err = id.MakeDirIfNecessary(sz.Hinweis, dir); err != nil {
			err = errors.Error(err)
			return
		}

		originalExt := sz.Stored.Zettel.AkteExt.String()
		originalFilename := filename
		filename = filename + ".md"

		inlineAkte := sz.Stored.Zettel.AkteExt.String() == "md"

		czs[i] = stored_zettel.CheckedOut{
			External: stored_zettel.External{
				Path: filename,
			},
		}

		c := zettel.FormatContextWrite{
			Zettel:            sz.Stored.Zettel,
			IncludeAkte:       inlineAkte,
			AkteReaderFactory: s.storeZettel,
		}

		if !inlineAkte && options.IncludeAkte {
			czs[i].External.AktePath = originalFilename + "." + originalExt
			c.ExternalAktePath = czs[i].External.AktePath
			c.IncludeAkte = true
		}

		if err = s.writeFormat(options, filename, c); err != nil {
			err = errors.Wrapped(err, "%s", sz.Named)
			stdprinter.Errf("%s (check out failed):\n", sz.Named)
			stdprinter.Error(err)
			continue
		}

		stdprinter.Outf("%s (checked out)\n", sz.Named)
	}

	return
}

func (s *Store) writeFormat(
	o CheckinOptions,
	p string,
	fc zettel.FormatContextWrite,
) (err error) {
	var f *os.File

	if f, err = open_file_guard.Create(p); err != nil {
		err = errors.Error(err)
		return
	}

	fc.Out = f

	defer open_file_guard.Close(f)

	if _, err = o.Format.WriteTo(fc); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
