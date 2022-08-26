package store_working_directory

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
	"github.com/friedenberg/zit/src/india/zettel_checked_out"
)

func (s *Store) Checkout(
	options CheckoutOptions,
	zs ...stored_zettel.Transacted,
) (czs []zettel_checked_out.CheckedOut, err error) {
	czs = make([]zettel_checked_out.CheckedOut, len(zs))

	for i, sz := range zs {
		if czs[i], err = s.CheckoutOne(options, sz); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}

func (s *Store) isInlineAkte(sz stored_zettel.Transacted) (isInline bool) {
	isInline = sz.Named.Stored.Zettel.TypOrDefault().String() == "md"

	if typKonfig, ok := s.Konfig.Typen[sz.Named.Stored.Zettel.Typ.String()]; ok {
		isInline = typKonfig.InlineAkte
	}

	return
}

func (s *Store) CheckoutOne(
	options CheckoutOptions,
	sz stored_zettel.Transacted,
) (cz zettel_checked_out.CheckedOut, err error) {
	var filename string

	if filename, err = id.MakeDirIfNecessary(sz.Named.Hinweis, s.cwd); err != nil {
		err = errors.Error(err)
		return
	}

	originalFilename := filename
	//TODO move user-fs-representation of zettel path to own function
	filename = filename + ".md"

	inlineAkte := s.isInlineAkte(sz)

	cz = zettel_checked_out.CheckedOut{
		External: stored_zettel.External{
			Path: filename,
		},
	}

	c := zettel.FormatContextWrite{
		Zettel:            sz.Named.Stored.Zettel,
		IncludeAkte:       inlineAkte,
		AkteReaderFactory: s.storeZettel,
	}

	switch options.CheckoutMode {
	case CheckoutModeAkteOnly:
		p := originalFilename + "." + sz.Named.Stored.Zettel.AkteExt()

		if err = s.writeAkte(sz.Named.Stored.Zettel.Akte, p); err != nil {
			err = errors.Error(err)
			return
		}

	case CheckoutModeZettelOnly:
		c.IncludeAkte = false

		fallthrough

	case CheckoutModeZettelAndAkte:
		c.IncludeAkte = true

		if !inlineAkte {
			cz.External.AktePath = originalFilename + "." + sz.Named.Stored.Zettel.AkteExt()
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
