package store_fs

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/echo/id"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/golf/typ"
	"github.com/friedenberg/zit/src/india/zettel"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/juliett/zettel_verzeichnisse"
)

type ZettelCheckedOutLogWriters struct {
	ZettelOnly collections.WriterFunc[*zettel_checked_out.Zettel]
	AkteOnly   collections.WriterFunc[*zettel_checked_out.Zettel]
	Both       collections.WriterFunc[*zettel_checked_out.Zettel]
}

func (s *Store) Checkout(
	options CheckoutOptions,
	ztw collections.WriterFunc[*zettel.Transacted],
) (zcs zettel_checked_out.MutableSet, err error) {
	zcs = zettel_checked_out.MakeMutableSetUnique(0)
	zts := zettel.MakeMutableSetUnique(0)

	if err = s.storeObjekten.Zettel().ReadAllSchwanzenVerzeichnisse(
		zettel_verzeichnisse.MakeWriterKonfig(s.konfig),
		zettel_verzeichnisse.MakeWriterZettelTransacted(ztw),
		zettel_verzeichnisse.MakeWriterZettelTransacted(zts.Add),
		zettel_verzeichnisse.MakeWriterZettelTransacted(
			collections.MakeWriterDoNotRepool[zettel.Transacted](),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	zts.Each(
		func(zt *zettel.Transacted) (err error) {
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
	case cz.Internal.Objekte.Equals(&cz.External.Objekte):
		cz.State = zettel_checked_out.StateJustCheckedOutButSame

	//TODO-P0 wait why?
	case cz.External.ZettelFD.Path == "":
		ok = true

	case options.Force || cz.State == zettel_checked_out.StateEmpty:
		ok = true
	}

	return
}

func (s Store) filenameForZettelTransacted(
	options CheckoutOptions,
	sz zettel.Transacted,
) (originalFilename string, filename string, err error) {
	if originalFilename, err = id.MakeDirIfNecessary(sz.Sku.Kennung, s.Cwd()); err != nil {
		err = errors.Wrap(err)
		return
	}

	filename = originalFilename + s.konfig.GetZettelFileExtension()

	return
}

func (s *Store) CheckoutOne(
	options CheckoutOptions,
	sz zettel.Transacted,
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
			//TODO-P2 handle fs state
			if err = s.zettelCheckedOutWriters.ZettelOnly(&cz); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	inlineAkte := typ.IsInlineAkte(sz.Objekte.Typ, s.konfig)

	cz = zettel_checked_out.Zettel{
		//TODO-P2 check diff with fs if already exists
		State:    zettel_checked_out.StateJustCheckedOut,
		Internal: sz,
		External: zettel_external.Zettel{
			ZettelFD: fd.FD{
				Path: filename,
			},
			Objekte: sz.Objekte,
			Sku: zettel_external.Sku{
				Sha:     sz.Sku.Sha,
				Kennung: sz.Sku.Kennung,
			},
		},
	}

	if !inlineAkte {
		t := sz.Objekte.Typ

		ty := s.konfig.GetTyp(t.String())

		if ty != nil {
			fe := ty.Objekte.Akte.FileExtension

			if fe == "" {
				fe = t.String()
			}

			cz.External.AkteFD = fd.FD{
				Path: originalFilename + "." + fe,
			}
		}
	}

	c := zettel.FormatContextWrite{
		Zettel:            sz.Objekte,
		IncludeAkte:       inlineAkte,
		AkteReaderFactory: s.storeObjekten,
	}

	switch options.CheckoutMode {
	case CheckoutModeAkteOnly:
		if err = s.writeAkte(
			cz.External.Objekte.Akte,
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
		cz.External.AkteFD = fd.FD{}

		if err = s.writeFormat(options, filename, c, &cz); err != nil {
			err = errors.Wrapf(err, "%s", sz.Sku)
			errors.PrintErrf("%s (check out failed):", sz.Sku)
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
