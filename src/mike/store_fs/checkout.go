package store_fs

import (
	"fmt"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/zettel_external"
	"github.com/friedenberg/zit/src/lima/zettel_checked_out"
)

func (s *Store) CheckoutQuery(
	options CheckoutOptions,
	ms kennung.MetaSet,
	f schnittstellen.FuncIter[objekte.CheckedOutLike],
) (err error) {
	if err = s.storeObjekten.Query(
		ms,
		func(t objekte.TransactedLike) (err error) {
			var co objekte.CheckedOutLike

			if co, err = s.checkoutOneGeneric(options, t); err != nil {
				err = errors.Wrap(err)
				return
			}

			return f(co)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) Checkout(
	options CheckoutOptions,
	ztw schnittstellen.FuncIter[*zettel.Transacted],
) (zcs zettel_checked_out.MutableSet, err error) {
	zcs = zettel_checked_out.MakeMutableSetUnique(0)
	zts := zettel.MakeMutableSetUnique(0)

	if err = s.storeObjekten.Zettel().ReadAllSchwanzen(
		iter.MakeChain(
			zettel.MakeWriterKonfig(s.erworben),
			ztw,
			collections.AddClone[zettel.Transacted, *zettel.Transacted](zts),
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
	case cz.Internal.Objekte.Equals(cz.External.Objekte):
		cz.State = objekte.CheckedOutStateJustCheckedOutButSame

	case options.Force || cz.State == objekte.CheckedOutStateEmpty:
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

	filename = originalFilename + s.erworben.GetZettelFileExtension()

	return
}

func (s *Store) checkoutOneGeneric(
	options CheckoutOptions,
	t objekte.TransactedLike,
) (co objekte.CheckedOutLike, err error) {
	switch tt := t.(type) {
	case *zettel.Transacted:
		return s.CheckoutOne(options, *tt)

	case *typ.Transacted:
		co, err = s.CheckoutOneTyp(options, *tt)
		s.checkedOutLogPrinter(co)

	default:
		err = gattung.MakeErrUnsupportedGattung(tt.GetSku2())
	}

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
		if cz, err = s.readOneFS(filename); err != nil {
			err = errors.Wrap(err)
			return
		}

		if !s.shouldCheckOut(options, cz) {
			// TODO-P2 handle fs state
			if err = s.checkedOutLogPrinter(&cz); err != nil {
				if errors.IsExist(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
					return
				}
			}

			return
		}
	}

	inlineAkte := s.erworben.IsInlineTyp(sz.Objekte.Typ)

	cz = zettel_checked_out.Zettel{
		// TODO-P2 check diff with fs if already exists
		CheckedOut: zettel.CheckedOut{
			State:    objekte.CheckedOutStateJustCheckedOut,
			Internal: sz,
			External: zettel.External{
				Objekte: sz.Objekte,
				Sku: zettel_external.Sku{
					ObjekteSha: sz.Sku.ObjekteSha,
					Kennung:    sz.Sku.Kennung,
				},
			},
		},
	}

	if options.CheckoutMode.IncludesZettel() {
		cz.External.FD.Path = filename
	}

	if !inlineAkte && options.CheckoutMode.IncludesAkte() {
		t := sz.Objekte.Typ

		ty := s.erworben.GetApproximatedTyp(t).ApproximatedOrActual()

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
		s.erworben,
	)

	if err = e.Encode(&cz.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.checkedOutLogPrinter(&cz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) CheckoutOneTyp(
	options CheckoutOptions,
	tk typ.Transacted,
) (co typ.CheckedOut, err error) {
	errors.TodoP1("extract format dependency")
	format := typ.MakeFormatText(s.storeObjekten)

	// var tt *typ.Transacted

	// if tt, err = s.storeObjekten.Typ().ReadOne(&tk); err != nil {
	// 	if errors.Is(err, objekte_store.ErrNotFound{}) {
	// 		err = nil
	// 		tt = &typ.Transacted{
	// 			Sku: sku.Transacted[kennung.Typ, *kennung.Typ]{
	// 				Kennung: tk,
	// 			},
	// 		}
	// 	} else {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}
	// }

	co.Internal = tk
	errors.TodoP0("external")

	var f *os.File

	if f, err = files.CreateExclusiveWriteOnly(
		path.Join(
			s.Cwd(),
			fmt.Sprintf("%s.%s", tk.Sku.Kennung, s.erworben.FileExtensions.Typ),
		),
	); err != nil {
		if errors.IsExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.DeferredCloser(&err, f)

	if co.External.FD, err = kennung.File(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = format.WriteFormat(f, &tk.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) MakeTempTypFiles(
	tks schnittstellen.Set[kennung.Typ],
) (ps []string, err error) {
	errors.TodoP3("add support for working directory")

	ps = make([]string, 0, tks.Len())

	format := typ.MakeFormatText(s.storeObjekten)

	if err = tks.Each(
		func(tk kennung.Typ) (err error) {
			var tt *typ.Transacted

			if tt, err = s.storeObjekten.Typ().ReadOne(&tk); err != nil {
				if errors.Is(err, objekte_store.ErrNotFound{}) {
					err = nil
					tt = &typ.Transacted{
						Sku: sku.Transacted[kennung.Typ, *kennung.Typ]{
							Kennung: tk,
						},
					}
				} else {
					err = errors.Wrap(err)
					return
				}
			}

			var f *os.File

			if f, err = files.CreateExclusiveWriteOnly(
				path.Join(
					s.Cwd(),
					fmt.Sprintf("%s.%s", tk.String(), s.erworben.FileExtensions.Typ),
				),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.Deferred(&err, f.Close)

			ps = append(ps, f.Name())

			if _, err = format.WriteFormat(f, &tt.Objekte); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
