package store_fs

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/foxtrot/id"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/zettel_checked_out"
)

func (s *Store) Checkout(
	options CheckoutOptions,
	ztw collections.WriterFunc[*zettel.Transacted],
) (zcs zettel_checked_out.MutableSet, err error) {
	zcs = zettel_checked_out.MakeMutableSetUnique(0)
	zts := zettel.MakeMutableSetUnique(0)

	if err = s.storeObjekten.Zettel().ReadAllSchwanzen(
		collections.MakeChain(
			zettel.MakeWriterKonfig(s.konfig),
			ztw,
			zts.AddAndDoNotRepool,
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = zts.Each(
		func(zt *zettel.Transacted) (err error) {
			var zc zettel_checked_out.Zettel

			if zc, err = s.CheckoutOne(options, *zt); err != nil {
				err = errors.Wrap(err)
				return
			}

			zcs.Add(&zc)
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Store) shouldCheckOut(
	options CheckoutOptions,
	cz zettel_checked_out.Zettel,
) (ok bool) {
	switch {
	case cz.Internal.Objekte.Equals(&cz.External.Objekte):
		cz.State = zettel_checked_out.StateJustCheckedOutButSame

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
			if err = s.zettelExternalLogPrinter(&cz.External); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	inlineAkte := s.konfig.IsInlineTyp(sz.Objekte.Typ)

	cz = zettel_checked_out.Zettel{
		//TODO-P2 check diff with fs if already exists
		State:    zettel_checked_out.StateJustCheckedOut,
		Internal: sz,
		External: zettel_external.Zettel{
			Objekte: sz.Objekte,
			Sku: zettel_external.Sku{
				ObjekteSha: sz.Sku.ObjekteSha,
				Kennung:    sz.Sku.Kennung,
			},
		},
	}

	if options.CheckoutMode.IncludesZettel() {
		cz.External.ZettelFD.Path = filename
	}

	if !inlineAkte && options.CheckoutMode.IncludesAkte() {
		t := sz.Objekte.Typ

		ty := s.konfig.GetTyp(t)

		var fe string

		if ty != nil {
			fe = ty.Objekte.Akte.FileExtension
		}

		if fe == "" {
			fe = t.String()
		}

		cz.External.AkteFD.Path = originalFilename + "." + fe
	}

	e := zettel_external.MakeFileEncoder(
		s.storeObjekten,
		s.konfig,
	)

	if err = e.Encode(&cz.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.zettelExternalLogPrinter(&cz.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
