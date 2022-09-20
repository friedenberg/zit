package store_working_directory

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/id"
	"github.com/friedenberg/zit/src/delta/id_set"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_external"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
)

func (s *Store) Checkout(
	options CheckoutOptions,
	ids id_set.Set,
) (zcs zettel_checked_out.Set, err error) {
	zcs = zettel_checked_out.MakeSetUnique(ids.Len())
	var zts zettel_transacted.Set

	if zts, err = s.storeObjekten.ReadMany(ids); err != nil {
		err = errors.Wrap(err)
		return
	}

	zts.Each(
		func(zt zettel_transacted.Zettel) (err error) {
			var zc zettel_checked_out.Zettel

			if zc, err = s.CheckoutOne(options, zt); err != nil {
				err = errors.Wrap(err)
				return
			}

			zcs.Add(zc)
			return
		},
	)

	return
}

func (s Store) shouldCheckOut(
	options CheckoutOptions,
	cz zettel_checked_out.Zettel,
) (ok bool) {
	if cz.External.ZettelFD.Path == "" {
		ok = true
		return
	}

	if cz.Internal.Named.Stored.Zettel.Equals(cz.External.Named.Stored.Zettel) {
		errors.Print(cz.Internal.Named.Stored.Zettel)
		errors.PrintOutf("%s (already checked out)", cz.Internal.Named)
		return
	}

	if options.Force || cz.State == zettel_checked_out.StateEmpty {
		ok = true
		return
	} else if cz.State == zettel_checked_out.StateExistsAndSame {
		errors.PrintOutf("%s (already checked out)", cz.Internal.Named)
	} else if cz.State == zettel_checked_out.StateExistsAndDifferent {
		errors.PrintOutf("%s (external has changes)", cz.Internal.Named)
	} else {
		errors.PrintOutf("%s (unknown state)", cz.Internal.Named)
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

	//TODO move user-fs-representation of zettel path to own function
	filename = originalFilename + ".md"

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
		},
	}

	c := zettel.FormatContextWrite{
		Zettel:            sz.Named.Stored.Zettel,
		IncludeAkte:       inlineAkte,
		AkteReaderFactory: s.storeObjekten,
	}

	cz.External.Named = sz.Named

	switch options.CheckoutMode {
	case CheckoutModeAkteOnly:
		p := originalFilename + "." + sz.Named.Stored.Zettel.AkteExt()
		cz.External.AkteFD.Path = p
		cz.External.ZettelFD.Path = ""

		if err = s.writeAkte(sz.Named.Stored.Zettel.Akte, p); err != nil {
			err = errors.Wrap(err)
			return
		}


	case CheckoutModeZettelAndAkte:
		c.IncludeAkte = true

		if !inlineAkte {
			cz.External.AkteFD.Path = originalFilename + "." + sz.Named.Stored.Zettel.AkteExt()
			c.ExternalAktePath = cz.External.AkteFD.Path
		}

		fallthrough

	case CheckoutModeZettelOnly:
		if err = s.writeFormat(options, filename, c); err != nil {
			err = errors.Wrapf(err, "%s", sz.Named)
			errors.PrintErrf("%s (check out failed):", sz.Named)
			return
		}

	default:
		err = errors.Errorf("unsupported checkout mode: %s", options.CheckoutMode)
		return
	}

	s.zettelCheckedOutPrinter.ZettelCheckedOut(cz).Print()

	return
}

func (s *Store) writeFormat(
	o CheckoutOptions,
	p string,
	fc zettel.FormatContextWrite,
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

	return
}

func (s *Store) writeAkte(
	sh sha.Sha,
	p string,
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

	errors.Print("starting io copy")
	if _, err = io.Copy(f, r); err != nil {
		errors.Print(" io copy faile")
		err = errors.Wrap(err)
		return
	}
	errors.Print(" io copy succeed")

	return
}
