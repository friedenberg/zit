package store_checkout

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
)

func (s *Store) Checkout(
	options CheckoutOptions,
	zs ...stored_zettel.Transacted,
) (czs []stored_zettel.CheckedOut, err error) {
	czs = make([]stored_zettel.CheckedOut, len(zs))

	for i, sz := range zs {
		if czs[i], err = s.CheckoutOne(options, sz); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}

func (s *Store) CheckoutOne(
	options CheckoutOptions,
	sz stored_zettel.Transacted,
) (cz stored_zettel.CheckedOut, err error) {
	var filename string

	if filename, err = id.MakeDirIfNecessary(sz.Hinweis, s.cwd); err != nil {
		err = errors.Error(err)
		return
	}

	originalExt := sz.Stored.Zettel.AkteExt.String()
	originalFilename := filename
	filename = filename + ".md"

	inlineAkte := sz.Stored.Zettel.AkteExt.String() == "md"

	if typKonfig, ok := s.Konfig.Typen[sz.Zettel.AkteExt.String()]; ok {
		inlineAkte = typKonfig.InlineAkte
	}

	cz = stored_zettel.CheckedOut{
		External: stored_zettel.External{
			Path: filename,
		},
	}

	c := zettel.FormatContextWrite{
		Zettel:            sz.Stored.Zettel,
		IncludeAkte:       inlineAkte,
		AkteReaderFactory: s.storeZettel,
	}

	switch options.CheckoutMode {
	case CheckoutModeAkteOnly:
		p := originalFilename + "." + originalExt

		if err = s.writeAkte(sz.Stored.Zettel.Akte, p); err != nil {
			err = errors.Error(err)
			return
		}

	case CheckoutModeZettelOnly:
		c.IncludeAkte = false

		fallthrough

	case CheckoutModeZettelAndAkte:
		c.IncludeAkte = true

		if !inlineAkte {
			cz.External.AktePath = originalFilename + "." + originalExt
			c.ExternalAktePath = cz.External.AktePath
			c.IncludeAkte = true
		}

		if err = s.writeFormat(options, filename, c); err != nil {
			err = errors.Wrapped(err, "%s", sz.Named)
			stdprinter.Errf("%s (check out failed):\n", sz.Named)
			stdprinter.Error(err)
			return
		}

	default:
		err = errors.Errorf("unsupported checkout mode: %s", options.CheckoutMode)
		return
	}

	stdprinter.Outf("%s (checked out)\n", sz.Named)

	return
}

func (s *Store) writeFormat(
	o CheckoutOptions,
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

func (s *Store) writeAkte(
	sh sha.Sha,
	p string,
) (err error) {
	var f *os.File

	if f, err = open_file_guard.Create(p); err != nil {
		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(f)

	var r io.ReadCloser

	if r, err = s.storeZettel.AkteReader(sh); err != nil {
		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(f)

	logz.Print("starting io copy")
	if _, err = io.Copy(f, r); err != nil {
		logz.Print(" io copy faile")
		err = errors.Error(err)
		return
	}
	logz.Print(" io copy succeed")

	return
}
