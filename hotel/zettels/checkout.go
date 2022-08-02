package zettels

import (
	"os"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

//TODO move to user_ops
func (zs *zettels) Checkout(options CheckinOptions, args ...string) (czs []stored_zettel.CheckedOut, err error) {
	var hins []hinweis.Hinweis
	var shas []sha.Sha

	if shas, hins, err = zs.hinweisen.ReadManyStrings(args...); err != nil {
		err = errors.Error(err)
		return
	}

	czs = make([]stored_zettel.CheckedOut, len(shas))

	var dir string

	if dir, err = os.Getwd(); err != nil {
		err = errors.Error(err)
		return
	}

	for i, sha := range shas {
		var sz stored_zettel.Named

		if sz, err = zs.Read(sha); err != nil {
			err = errors.Error(err)
			return
		}

		var filename string

		if filename, err = id.MakeDirIfNecessary(hins[i], dir); err != nil {
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
			AkteReaderFactory: zs,
		}

		if !inlineAkte && options.IncludeAkte {
			czs[i].External.AktePath = originalFilename + "." + originalExt
			c.ExternalAktePath = czs[i].External.AktePath
			c.IncludeAkte = true
		}

		if err = zs.writeFormat(options, filename, c); err != nil {
			err = errors.Errorf("%s: %s", sz.Hinweis, err)
      stdprinter.Errf("[%s %s] (check out failed):\n", hins[i], shas[i])
			stdprinter.Error(err)
			continue
		}

		stdprinter.Outf("[%s %s] (checked out)\n", hins[i], shas[i])
	}

	return
}

func (zs zettels) writeFormat(o CheckinOptions, p string, fc zettel.FormatContextWrite) (err error) {
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
