package store_working_directory

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
)

func (s *Store) Checkout(
	options CheckoutOptions,
	zs ...zettel_transacted.Zettel,
) (czs []zettel_checked_out.Zettel, err error) {
	czs = make([]zettel_checked_out.Zettel, len(zs))

	for i, sz := range zs {
		if czs[i], err = s.CheckoutOne(options, sz); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) isInlineAkte(sz zettel_transacted.Zettel) (isInline bool) {
	isInline = sz.Named.Stored.Zettel.TypOrDefault().String() == "md"

	if typKonfig, ok := s.Konfig.Typen[sz.Named.Stored.Zettel.Typ.String()]; ok {
		isInline = typKonfig.InlineAkte
	}

	return
}

func (s *Store) CheckoutOne(
	options CheckoutOptions,
	sz zettel_transacted.Zettel,
) (cz zettel_checked_out.Zettel, err error) {
	var filename string

	if filename, err = id.MakeDirIfNecessary(sz.Named.Hinweis, s.cwd); err != nil {
		err = errors.Wrap(err)
		return
	}

	originalFilename := filename
	//TODO move user-fs-representation of zettel path to own function
	filename = filename + ".md"

	inlineAkte := s.isInlineAkte(sz)

	cz = zettel_checked_out.Zettel{
		External: zettel_external.Zettel{
			ZettelFD: zettel_external.FD{
				Path: filename,
			},
		},
	}

	c := zettel.FormatContextWrite{
		Zettel:            sz.Named.Stored.Zettel,
		IncludeAkte:       inlineAkte,
		AkteReaderFactory: s.storeObjekten,
	}

	switch options.CheckoutMode {
	case CheckoutModeAkteOnly:
		p := originalFilename + "." + sz.Named.Stored.Zettel.AkteExt()

		if err = s.writeAkte(sz.Named.Stored.Zettel.Akte, p); err != nil {
			err = errors.Wrap(err)
			return
		}

	case CheckoutModeZettelOnly:
		c.IncludeAkte = false

		fallthrough

	case CheckoutModeZettelAndAkte:
		c.IncludeAkte = true

		if !inlineAkte {
			cz.External.AkteFD.Path = originalFilename + "." + sz.Named.Stored.Zettel.AkteExt()
			c.ExternalAktePath = cz.External.AkteFD.Path
			c.IncludeAkte = true
		}

		if err = s.writeFormat(options, filename, c); err != nil {
			err = errors.Wrapf(err, "%s", sz.Named)
			errors.PrintErrf("%s (check out failed):", sz.Named)
			errors.PrintErr(err)
			return
		}

	default:
		err = errors.Errorf("unsupported checkout mode: %s", options.CheckoutMode)
		return
	}

	errors.PrintOutf("%s (checked out)", sz.Named)

	return
}

func (s *Store) writeFormat(
	o CheckoutOptions,
	p string,
	fc zettel.FormatContextWrite,
) (err error) {
	var f *os.File

	if f, err = open_file_guard.Create(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	fc.Out = f

	defer open_file_guard.Close(f)

	if _, err = o.Format.WriteTo(fc); err != nil {
		err = errors.Wrap(err)
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
		err = errors.Wrap(err)
		return
	}

	defer open_file_guard.Close(f)

	var r io.ReadCloser

	if r, err = s.storeObjekten.AkteReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer open_file_guard.Close(f)

	errors.Print("starting io copy")
	if _, err = io.Copy(f, r); err != nil {
		errors.Print(" io copy faile")
		err = errors.Wrap(err)
		return
	}
	errors.Print(" io copy succeed")

	return
}
