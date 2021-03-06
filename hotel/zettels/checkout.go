package zettels

import (
	"os"
)

//TODO move to user_ops
func (zs *zettels) Checkout(options CheckinOptions, args ...string) (czs []_ZettelCheckedOut, err error) {
	var hins []_Hinweis
	var shas []_Sha

	if shas, hins, err = zs.hinweisen.ReadManyStrings(args...); err != nil {
		err = _Error(err)
		return
	}

	czs = make([]_ZettelCheckedOut, len(shas))

	var dir string

	if dir, err = os.Getwd(); err != nil {
		err = _Error(err)
		return
	}

	for i, sha := range shas {
		var sz _NamedZettel

		if sz, err = zs.Read(sha); err != nil {
			err = _Error(err)
			return
		}

		var filename string

		if filename, err = _IdMakeDirNecessary(hins[i], dir); err != nil {
			err = _Error(err)
			return
		}

		originalExt := sz.Stored.Zettel.AkteExt.String()
		originalFilename := filename
		filename = filename + ".md"

		inlineAkte := sz.Stored.Zettel.AkteExt.String() == "md"

		czs[i] = _ZettelCheckedOut{
			External: _ZettelExternal{
				Path: filename,
			},
		}

		c := _ZettelFormatContextWrite{
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
			err = _Errorf("%s: %s", sz.Hinweis, err)
			_Errf("[%s %s] (check out failed)\n", hins[i], shas[i], err)
			continue
		}

		_Outf("[%s %s] (checked out)\n", hins[i], shas[i])
	}

	return
}

func (zs zettels) writeFormat(o CheckinOptions, p string, fc _ZettelFormatContextWrite) (err error) {
	var f *os.File

	if f, err = _Create(p); err != nil {
		err = _Error(err)
		return
	}

	fc.Out = f

	defer _Close(f)

	if _, err = o.Format.WriteTo(fc); err != nil {
		err = _Error(err)
		return
	}

	return
}
