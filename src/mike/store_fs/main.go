package store_fs

import (
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/kasten"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/lima/store_objekten"
)

type Store struct {
	sonnenaufgang kennung.Time
	erworben      konfig.Compiled
	standort.Standort

	storeObjekten *store_objekten.Store

	checkedOutLogPrinter schnittstellen.FuncIter[objekte.CheckedOutLike]
}

func New(
	t kennung.Time,
	k konfig.Compiled,
	st standort.Standort,
	storeObjekten *store_objekten.Store,
) (s *Store, err error) {
	s = &Store{
		sonnenaufgang: t,
		erworben:      k,
		Standort:      st,
		storeObjekten: storeObjekten,
	}

	return
}

func (s *Store) SetCheckedOutLogPrinter(
	zelw schnittstellen.FuncIter[objekte.CheckedOutLike],
) {
	s.checkedOutLogPrinter = zelw
}

// TODO-P3 move to standort
func (s Store) IndexFilePath() string {
	return path.Join(s.Cwd(), ".ZitCheckoutStoreIndex")
}

func (s Store) Flush() (err error) {
	return
}

// Methods
// ReadOne
// ReadMany
// ReadManyHistory
// func (s *Store) ReadOne(h kennung.Hinweis) (zt *zettel.Transacted, err error) {
// 	errors.TodoP1("include cwd sigil")
// 	if zt, err = s.storeObjekten.Zettel().ReadOne(&h); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if !s.erworben.IncludeCwd {
// 		return
// 	}

// 	var pz cwd.CwdFiles

// 	if pz, err = cwd.MakeCwdFilesExactly(
// 		s.erworben,
// 		s.Standort.Cwd(),
// 		fmt.Sprintf("%s.%s", h, s.erworben.FileExtensions.Zettel),
// 	); err != nil {
// 		if errors.IsNotExist(err) {
// 			err = nil
// 		} else {
// 			err = errors.Wrap(err)
// 		}

// 		return
// 	}

// 	var checked_out zettel.CheckedOut

// 	var readFunc func() (zettel.External, error)

// 	p := pz.Zettelen.Any()

// 	switch {
// 	case p.GetAkteFD().Path == "":
// 		readFunc = func() (zettel.External, error) {
// 			return s.storeObjekten.Zettel().ReadOneExternal.Read(p.GetObjekteFD().Path)
// 		}

// 	case p.GetObjekteFD().Path == "":
// 		readFunc = func() (zettel.CheckedOut, error) {
// 			return s.ReadExternalZettelFromAktePath(p.GetAkteFD().Path)
// 		}

// 	default:
// 		// TODO-P3 validate that the zettel file points to the akte in the metadatei
// 		readFunc = func() (zettel.CheckedOut, error) {
// 			return s.Read(p.GetObjekteFD().Path)
// 		}
// 	}

// 	if checked_out, err = readFunc(); err != nil {
// 		if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
// 			err = nil
// 		} else {
// 			err = errors.Wrap(err)
// 		}

// 		return
// 	}

// 	zt.Sku = checked_out.External.Sku.Transacted()
// 	zt.Objekte = checked_out.External.Objekte
// 	zt.Sku.Schwanz = s.sonnenaufgang

// 	return
// }

func (s *Store) ReadFiles(
	fs *cwd.CwdFiles,
	ms kennung.MetaSet,
	f schnittstellen.FuncIter[objekte.CheckedOutLike],
) (err error) {
	zettelEMGR := objekte_store.MakeExternalMaybeGetterReader[
		zettel.Objekte,
		*zettel.Objekte,
		kennung.Hinweis,
		*kennung.Hinweis,
		zettel.Verzeichnisse,
		*zettel.Verzeichnisse,
	](
		fs.GetZettel,
		s.storeObjekten.Zettel(),
	)

	etikettEMGR := objekte_store.MakeExternalMaybeGetterReader[
		etikett.Akte,
		*etikett.Akte,
		kennung.Etikett,
		*kennung.Etikett,
		objekte.NilVerzeichnisse[etikett.Akte],
		*objekte.NilVerzeichnisse[etikett.Akte],
	](
		fs.GetEtikett,
		s.storeObjekten.Etikett(),
	)

	typEMGR := objekte_store.MakeExternalMaybeGetterReader[
		typ.Akte,
		*typ.Akte,
		kennung.Typ,
		*kennung.Typ,
		objekte.NilVerzeichnisse[typ.Akte],
		*objekte.NilVerzeichnisse[typ.Akte],
	](
		fs.GetTyp,
		s.storeObjekten.Typ(),
	)

	kastenEMGR := objekte_store.MakeExternalMaybeGetterReader[
		kasten.Akte,
		*kasten.Akte,
		kennung.Kasten,
		*kennung.Kasten,
		kasten.Verzeichnisse,
		*kasten.Verzeichnisse,
	](
		fs.GetKasten,
		s.storeObjekten.Kasten(),
	)

	if err = s.storeObjekten.Query(
		ms,
		iter.MakeChain(
			func(e objekte.TransactedLikePtr) (err error) {
				var col objekte.CheckedOutLikePtr

				switch et := e.(type) {
				case *zettel.Transacted:
					if col, err = zettelEMGR.ReadOne(*et); err != nil {
						var errAkte store_objekten.ErrExternalAkteExtensionMismatch

						if errors.As(err, &errAkte) {
							fs.MarkUnsureAkten(errAkte.Actual)
							err = nil
						} else {
							err = errors.Wrap(err)
						}

						return
					}

				case *typ.Transacted:
					if col, err = typEMGR.ReadOne(*et); err != nil {
						err = errors.Wrap(err)
						return
					}

				case *kasten.Transacted:
					if col, err = kastenEMGR.ReadOne(*et); err != nil {
						err = errors.Wrap(err)
						return
					}

				case *etikett.Transacted:
					if col, err = etikettEMGR.ReadOne(*et); err != nil {
						err = errors.Wrap(err)
						return
					}

				case *erworben.Transacted:
					errors.TodoP1("implement checked out konfig?")
					return

				default:
					err = errors.Implement()
					return
				}

				col.DetermineState()

				if err = f(col); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = fs.EachCreatableMatchable(
		iter.MakeChain(
			func(ilg kennung.IdLikeGetter) (err error) {
				switch il := ilg.(type) {
				case cwd.Kasten:
					if err = s.storeObjekten.GetAbbrStore().KastenExists(
						il.Kennung,
					); err == nil {
						err = iter.MakeErrStopIteration()
						return
					}

					err = nil

					var tco kasten.CheckedOut

					if tco.External, err = s.storeObjekten.Kasten().ReadOneExternal(
						il,
						nil,
					); err != nil {
						if errors.IsNotExist(err) {
							err = iter.MakeErrStopIteration()
						} else {
							err = errors.Wrapf(err, "CwdEtikett: %#v", il)
						}

						return
					}

					tco.State = objekte.CheckedOutStateUntracked

					if err = f(&tco); err != nil {
						err = errors.Wrap(err)
						return
					}

				case cwd.Typ:
					if err = s.storeObjekten.GetAbbrStore().TypExists(
						il.Kennung,
					); err == nil {
						err = iter.MakeErrStopIteration()
						return
					}

					err = nil

					var tco typ.CheckedOut

					if tco.External, err = s.storeObjekten.Typ().ReadOneExternal(
						il,
						nil,
					); err != nil {
						if errors.IsNotExist(err) {
							err = iter.MakeErrStopIteration()
						} else {
							err = errors.Wrapf(err, "CwdEtikett: %#v", il)
						}

						return
					}

					tco.State = objekte.CheckedOutStateUntracked

					if err = f(&tco); err != nil {
						err = errors.Wrap(err)
						return
					}

				case cwd.Etikett:
					if err = s.storeObjekten.GetAbbrStore().EtikettExists(
						il.Kennung,
					); err == nil {
						err = iter.MakeErrStopIteration()
						return
					}

					err = nil

					var tco etikett.CheckedOut

					if tco.External, err = s.storeObjekten.Etikett().ReadOneExternal(
						il,
						nil,
					); err != nil {
						if errors.IsNotExist(err) {
							err = iter.MakeErrStopIteration()
						} else {
							err = errors.Wrapf(err, "CwdEtikett: %#v", il)
						}

						return
					}

					tco.State = objekte.CheckedOutStateUntracked

					if err = f(&tco); err != nil {
						err = errors.Wrap(err)
						return
					}

				default:
					err = errors.Errorf("unsupported id like: %T", il)
				}

				return
			},
			// func(ilg sku.IdLikeGetter) (err error) {
			// },
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// if cz.State > objekte.CheckedOutStateExistsAndSame {
// TODO-P4 rewrite with verzeichnisseAll
// exSha := cz.External.Sku.Sha
// cz.Matches.Zettelen, _ = s.storeObjekten.ReadZettelSha(exSha)
// cz.Matches.Zettelen, _ = cz.Matches.Zettelen.Filter(nil, filter)

// exAkteSha := cz.External.Objekte.Akte
// cz.Matches.Akten, _ = s.storeObjekten.ReadAkteSha(exAkteSha)
// cz.Matches.Akten, _ = cz.Matches.Akten.Filter(nil, filter)

// bez := cz.External.Objekte.Bezeichnung.String()
// cz.Matches.Bezeichnungen, _ = s.storeObjekten.ReadBezeichnung(bez)
// cz.Matches.Bezeichnungen, _ = cz.Matches.Bezeichnungen.Filter(nil, filter)
// }
