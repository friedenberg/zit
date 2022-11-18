package store_fs

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
)

type ZettelCheckedOutLogWriters struct {
	ZettelOnly collections.WriterFunc[*zettel_checked_out.Zettel]
	AkteOnly   collections.WriterFunc[*zettel_checked_out.Zettel]
	Both       collections.WriterFunc[*zettel_checked_out.Zettel]
}

func (s *Store) Checkout(
	options CheckoutOptions,
	ztw collections.WriterFunc[*zettel_transacted.Zettel],
) (zcs zettel_checked_out.MutableSet, err error) {
	zcs = zettel_checked_out.MakeMutableSetUnique(0)
	zts := zettel_transacted.MakeMutableSetUnique(0)

	if err = s.storeObjekten.ReadAllSchwanzenTransacted(
		ztw,
		zts.Add,
		collections.MakeWriterDoNotRepool[zettel_transacted.Zettel](),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	zts.Each(
		func(zt *zettel_transacted.Zettel) (err error) {
			var zc zettel_checked_out.Zettel

			if zc, err = s.CheckoutOne(options, *zt); err != nil {
				err = errors.Wrap(err)
				return
			}

			zcs.Add(&zc)
			return
		},
	)

	return
}

func (s Store) shouldCheckOut(
	options CheckoutOptions,
	cz zettel_checked_out.Zettel,
) (ok bool) {

	switch {
	case cz.Internal.Named.Stored.Zettel.Equals(cz.External.Named.Stored.Zettel):
		cz.State = zettel_checked_out.StateJustCheckedOutButSame

	//TODO wait why?
	case cz.External.ZettelFD.Path == "":
		ok = true

	case options.Force || cz.State == zettel_checked_out.StateEmpty:
		ok = true
	}

	return
}

func (s Store) filenameForZettelTransacted(
	options CheckoutOptions,
	sz zettel_transacted.Zettel,
) (originalFilename string, filename string, err error) {
	if originalFilename, err = id.MakeDirIfNecessary(sz.Named.Hinweis, s.cwd); err != nil {
		err = errors.Wrap(err)
		return
	}

	filename = originalFilename + s.Konfig.Compiled.GetZettelFileExtension()

	return
}

func (s *Store) CheckoutOne(
	options CheckoutOptions,
	sz zettel_transacted.Zettel,
) (cz zettel_checked_out.Zettel, err error) {
	var originalFilename, filename string

	if originalFilename, filename, err = s.filenameForZettelTransacted(options, sz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if files.Exists(filename) {
		if cz, err = s.Read(filename); err != nil {
			err = errors.Wrap(err)
			return
		}

		if !s.shouldCheckOut(options, cz) {
			//TODO handle fs state
			if err = s.zettelCheckedOutWriters.ZettelOnly(&cz); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	inlineAkte := sz.Named.Stored.Zettel.IsInlineAkte(s.Konfig.Konfig)

	cz = zettel_checked_out.Zettel{
		//TODO check diff with fs if already exists
		State:    zettel_checked_out.StateJustCheckedOut,
		Internal: sz,
		External: zettel_external.Zettel{
			ZettelFD: zettel_external.FD{
				Path: filename,
			},
			Named: sz.Named,
		},
		Matches: zettel_checked_out.MakeMatches(),
	}

	if !inlineAkte {
		t := sz.Named.Stored.Zettel.Typ

		ty := s.Compiled.GetType(t.String())

		if ty != nil && ty.FileExtension != "" {
			cz.External.AkteFD = zettel_external.FD{
				Path: originalFilename + "." + ty.FileExtension,
			}
		}
	}

	c := zettel.FormatContextWrite{
		Zettel:            sz.Named.Stored.Zettel,
		IncludeAkte:       inlineAkte,
		AkteReaderFactory: s.storeObjekten,
	}

	switch options.CheckoutMode {
	case CheckoutModeAkteOnly:
		if err = s.writeAkte(
			cz.External.Named.Stored.Zettel.Akte,
			cz.External.AkteFD.Path,
			&cz,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	case CheckoutModeZettelAndAkte:
		c.IncludeAkte = true

		if !inlineAkte {
			c.ExternalAktePath = cz.External.AkteFD.Path
		}

		fallthrough

	case CheckoutModeZettelOnly:
		if err = s.writeFormat(options, filename, c, &cz); err != nil {
			err = errors.Wrapf(err, "%s", sz.Named)
			errors.PrintErrf("%s (check out failed):", sz.Named)
			return
		}

	default:
		err = errors.Errorf("unsupported checkout mode: %s", options.CheckoutMode)
		return
	}

	return
}

func (s *Store) writeFormat(
	o CheckoutOptions,
	p string,
	fc zettel.FormatContextWrite,
	zco *zettel_checked_out.Zettel,
) (err error) {
	var f *os.File

	if f, err = files.Create(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	fc.Out = f

	defer files.Close(f)

	if _, err = o.Format.WriteTo(fc); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fc.ExternalAktePath == "" {
		if err = s.zettelCheckedOutWriters.ZettelOnly(zco); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = s.zettelCheckedOutWriters.Both(zco); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) writeAkte(
	sh sha.Sha,
	p string,
	zco *zettel_checked_out.Zettel,
) (err error) {
	var f *os.File

	if f, err = files.Create(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer files.Close(f)

	var r io.ReadCloser

	if r, err = s.storeObjekten.AkteReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer files.Close(f)

	if _, err = io.Copy(f, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.zettelCheckedOutWriters.AkteOnly(zco); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
