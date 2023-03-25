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
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/kilo/zettel_external"
	"github.com/friedenberg/zit/src/lima/store_objekten"
)

func (s *Store) CheckoutQuery(
	options CheckoutOptions,
	ms kennung.MetaSet,
	f schnittstellen.FuncIter[objekte.CheckedOutLike],
) (err error) {
	if err = s.storeObjekten.Query(
		ms,
		func(t objekte.TransactedLike) (err error) {
			var co objekte.CheckedOutLikePtr

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
) (zcs zettel.MutableSetCheckedOut, err error) {
	zcs = zettel.MakeMutableSetCheckedOutUnique(0)
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
			var zc zettel.CheckedOut

			if zc, err = s.CheckoutOneZettel(options, *zt); err != nil {
				err = errors.Wrap(err)
				return
			}

			zcs.Add(zc)
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
	cz zettel.CheckedOut,
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
) (cop objekte.CheckedOutLikePtr, err error) {
	switch tt := t.(type) {
	case *zettel.Transacted:
		var co zettel.CheckedOut
		co, err = s.CheckoutOneZettel(options, *tt)
		cop = &co

	case *typ.Transacted:
		cop, err = s.storeObjekten.Typ().CheckoutOne(store_objekten.CheckoutOptions(options), tt)

	case *etikett.Transacted:
		cop, err = s.storeObjekten.Etikett().CheckoutOne(store_objekten.CheckoutOptions(options), tt)

	default:
		// err = errors.Implement()
		err = gattung.MakeErrUnsupportedGattung(tt.GetSku2())
		return
	}

	cop.DetermineState()
	s.checkedOutLogPrinter(cop)

	return
}

func (s *Store) CheckoutOneZettel(
	options CheckoutOptions,
	sz zettel.Transacted,
) (cz zettel.CheckedOut, err error) {
	cz.Internal = sz

	var originalFilename, filename string

	if originalFilename, filename, err = s.filenameForZettelTransacted(options, sz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if files.Exists(filename) {
		var e cwd.Zettel
		ok := false

		if e, ok = options.Cwd.GetZettel(sz.Sku.Kennung); !ok {
			err = errors.Errorf("file at %s not recognized as zettel: %s", filename, sz)
			return
		}

		if cz.External, err = s.storeObjekten.Zettel().ReadOneExternal(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		cz.DetermineState()

		if !s.shouldCheckOut(options, cz) {
			return
		}
	}

	inlineAkte := s.erworben.IsInlineTyp(sz.Objekte.Typ)

	cz.State = objekte.CheckedOutStateJustCheckedOut
	cz.External = zettel.External{
		Objekte: sz.Objekte,
		Sku: zettel_external.Sku{
			ObjekteSha: sz.Sku.ObjekteSha,
			Kennung:    sz.Sku.Kennung,
		},
	}

	if options.CheckoutMode.IncludesObjekte() {
		cz.External.Sku.FDs.Objekte.Path = filename
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

		cz.External.Sku.FDs.Akte.Path = originalFilename + "." + fe
	}

	e := zettel_external.MakeFileEncoder(
		s.storeObjekten,
		s.erworben,
	)

	if err = e.Encode(&cz.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// func (s *Store) CheckoutOneEtikett(
// 	options CheckoutOptions,
// 	tk etikett.Transacted,
// ) (co etikett.CheckedOut, err error) {
// 	errors.TodoP1("extract format dependency")
// 	format := etikett.MakeFormatText(s.storeObjekten)

// 	co.Internal = tk
// 	co.External.Sku = tk.Sku.GetExternal()

// 	var f *os.File

// 	p := path.Join(
// 		s.Cwd(),
// 		fmt.Sprintf("%s.%s", tk.Sku.Kennung, s.erworben.FileExtensions.Etikett),
// 	)

// 	if f, err = files.CreateExclusiveWriteOnly(p); err != nil {
// 		if errors.IsExist(err) {
// 			if co.External, err = s.storeObjekten.Etikett().ReadOneExternal(
// 				cwd.Etikett{
// 					Kennung: tk.Sku.Kennung,
// 					FDs: sku.ExternalFDs{
// 						Objekte: kennung.FD{
// 							Path: p,
// 						},
// 					},
// 				},
// 			); err != nil {
// 				err = errors.Wrap(err)
// 				return
// 			}

// 			co.External.Sku.Kennung = tk.Sku.Kennung
// 		} else {
// 			err = errors.Wrap(err)
// 		}

// 		return
// 	}

// 	defer errors.DeferredCloser(&err, f)

// 	if co.External.Sku.FDs.Objekte, err = kennung.File(f); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if _, err = format.Format(f, &tk.Objekte); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

// func (s *Store) CheckoutOneTyp(
// 	options CheckoutOptions,
// 	tk typ.Transacted,
// ) (co typ.CheckedOut, err error) {
// 	errors.TodoP1("extract format dependency")
// 	format := typ.MakeFormatText(s.storeObjekten)

// 	co.Internal = tk
// 	co.External.Sku = tk.Sku.GetExternal()

// 	var f *os.File

// 	p := path.Join(
// 		s.Cwd(),
// 		fmt.Sprintf("%s.%s", tk.Sku.Kennung, s.erworben.FileExtensions.Typ),
// 	)

// 	if f, err = files.CreateExclusiveWriteOnly(p); err != nil {
// 		if errors.IsExist(err) {
// 			if co.External, err = s.storeObjekten.Typ().ReadOneExternal(
// 				cwd.Typ{
// 					Kennung: tk.Sku.Kennung,
// 					FDs: sku.ExternalFDs{
// 						Objekte: kennung.FD{
// 							Path: p,
// 						},
// 					},
// 				},
// 			); err != nil {
// 				err = errors.Wrap(err)
// 				return
// 			}

// 			co.External.Sku.Kennung = tk.Sku.Kennung
// 		} else {
// 			err = errors.Wrap(err)
// 		}

// 		return
// 	}

// 	defer errors.DeferredCloser(&err, f)

// 	if co.External.Sku.FDs.Objekte, err = kennung.File(f); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if _, err = format.Format(f, &tk.Objekte); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

// TODO discard
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

			if _, err = format.Format(f, &tt.Objekte); err != nil {
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
