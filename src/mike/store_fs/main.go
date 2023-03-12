package store_fs

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/hinweisen"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/kilo/zettel_external"
	"github.com/friedenberg/zit/src/lima/store_objekten"
	"github.com/friedenberg/zit/src/todo"
)

type Store struct {
	sonnenaufgang ts.Time
	erworben      konfig.Compiled
	standort.Standort

	format zettel.ObjekteFormat

	storeObjekten *store_objekten.Store

	zettelExternalLogPrinter schnittstellen.FuncIter[*zettel_external.Zettel]
	checkedOutLogPrinter     schnittstellen.FuncIter[objekte.CheckedOutLike]

	entries      map[string]Entry
	indexWasRead bool
	hasChanges   bool
}

func New(
	t ts.Time,
	k konfig.Compiled,
	st standort.Standort,
	storeObjekten *store_objekten.Store,
) (s *Store, err error) {
	s = &Store{
		sonnenaufgang: t,
		erworben:      k,
		Standort:      st,
		format: zettel.MakeObjekteTextFormat(
			storeObjekten,
			nil,
		),
		storeObjekten: storeObjekten,
		entries:       make(map[string]Entry),
	}

	return
}

func (s *Store) SetCheckedOutLogPrinter(
	zelw schnittstellen.FuncIter[objekte.CheckedOutLike],
) {
	s.checkedOutLogPrinter = zelw
}

func (s *Store) SetZettelExternalLogPrinter(
	zelw schnittstellen.FuncIter[*zettel_external.Zettel],
) {
	s.zettelExternalLogPrinter = zelw
}

// TODO-P3 move to standort
func (s Store) IndexFilePath() string {
	return path.Join(s.Cwd(), ".ZitCheckoutStoreIndex")
}

func (s Store) flushToTemp() (tfp string, err error) {
	var f *os.File

	if f, err = files.TempFile(s.Standort.DirTempLocal()); err != nil {
		err = errors.Wrap(err)
		return
	}

	tfp = f.Name()

	defer errors.Deferred(&err, f.Close)

	w := bufio.NewWriter(f)
	defer errors.Deferred(&err, w.Flush)

	for p, e := range s.entries {
		out := fmt.Sprintf("%s %s\n", p, e)
		errors.Log().Printf("flushing zettel: %q", out)
		w.WriteString(fmt.Sprint(out))
	}

	return
}

func (s Store) Flush() (err error) {
	if s.hasChanges {
		var tfp string

		if tfp, err = s.flushToTemp(); err != nil {
			err = errors.Wrap(err)
			return
		}

		errors.Log().Printf("renaming %s to %s", tfp, s.IndexFilePath())
		if err = os.Rename(tfp, s.IndexFilePath()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s Store) MakeExternalZettelFromZettel(
	p string,
) (ez zettel.External, err error) {
	if p, err = filepath.Abs(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if p, err = filepath.Rel(s.Cwd(), p); err != nil {
		err = errors.Wrap(err)
		return
	}

	ez.Sku.FDs.Objekte.Path = p

	head, tail := id.HeadTailFromFileName(p)

	if ez.Sku.Kennung, err = kennung.MakeHinweis(head + "/" + tail); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Store) readZettelFromFile(ez *zettel.External) (err error) {
	if !files.Exists(ez.GetObjekteFD().Path) {
		// if the path does not have an extension, try looking for a file with that
		// extension
		// TODO-P4 modify this to use globs
		if filepath.Ext(ez.GetObjekteFD().Path) == "" {
			ez.Sku.FDs.Objekte.Path = ez.GetObjekteFD().Path + s.erworben.GetZettelFileExtension()
			return s.readZettelFromFile(ez)
		}

		err = os.ErrNotExist

		return
	}

	c := zettel.ObjekteParserContext{}

	var f *os.File

	if f, err = files.Open(ez.GetObjekteFD().Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	if _, err = s.format.Parse(f, &c); err != nil {
		err = errors.Wrapf(err, "%s", f.Name())
		return
	}

	if ez.Sku.ObjekteSha, err = s.storeObjekten.Zettel().WriteZettelObjekte(
		c.Zettel,
	); err != nil {
		err = errors.Wrapf(err, "%s", f.Name())
		return
	}

	ez.Objekte = c.Zettel
	ez.Sku.FDs.Akte.Path = c.AktePath

	unrecoverableErrors := errors.MakeMulti()

	for _, e := range errors.Split(c.Errors) {
		// var errAkteInlineAndFilePath zettel.ErrHasInlineAkteAndFilePath

		// if errors.As(e, &errAkteInlineAndFilePath) {
		// 	var z1 zettel.Zettel

		// 	if z1, err = errAkteInlineAndFilePath.Recover(); err != nil {
		// 		unrecoverableErrors.Add(errors.Wrap(err))
		// 		continue
		// 	}

		// 	ez.Objekte = z1
		// 	continue
		// }

		var err1 zettel.ErrHasInvalidAkteShaOrFilePath

		if errors.As(e, &err1) {
			var mutter *zettel.Transacted

			if mutter, err = s.storeObjekten.Zettel().ReadOne(ez.Sku.Kennung); err != nil {
				unrecoverableErrors.Add(errors.Wrap(err))
				continue
			}

			ez.Objekte.Akte = mutter.Objekte.Akte

			continue
		}

		unrecoverableErrors.Add(e)
	}

	if !unrecoverableErrors.Empty() {
		err = errors.Wrap(unrecoverableErrors)
		return
	}

	return
}

// Methods
// ReadOne
// ReadMany
// ReadManyHistory
func (s *Store) ReadOne(h kennung.Hinweis) (zt *zettel.Transacted, err error) {
	errors.TodoP0("include cwd sigil")
	if zt, err = s.storeObjekten.Zettel().ReadOne(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !s.erworben.IncludeCwd {
		return
	}

	var pz cwd.CwdFiles

	if pz, err = cwd.MakeCwdFilesExactly(
		s.erworben,
		s.Standort.Cwd(),
		fmt.Sprintf("%s.%s", h, s.erworben.FileExtensions.Zettel),
	); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	var checked_out zettel.CheckedOut

	var readFunc func() (zettel.CheckedOut, error)

	p := pz.Zettelen.Any()

	switch {
	case p.GetAkteFD().Path == "":
		readFunc = func() (zettel.CheckedOut, error) {
			return s.Read(p.GetObjekteFD().Path)
		}

	case p.GetObjekteFD().Path == "":
		readFunc = func() (zettel.CheckedOut, error) {
			return s.ReadExternalZettelFromAktePath(p.GetAkteFD().Path)
		}

	default:
		// TODO-P3 validate that the zettel file points to the akte in the metadatei
		readFunc = func() (zettel.CheckedOut, error) {
			return s.Read(p.GetObjekteFD().Path)
		}
	}

	if checked_out, err = readFunc(); err != nil {
		if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	zt.Sku = checked_out.External.Sku.Transacted()
	zt.Objekte = checked_out.External.Objekte
	zt.Sku.Schwanz = s.sonnenaufgang

	return
}

func (s *Store) ReadMany(
	w1 schnittstellen.FuncIter[*zettel.Transacted],
) (err error) {
	w := w1

	if s.erworben.IncludeCwd {
		w = func(z *zettel.Transacted) (err error) {
			// TODO-P2 akte fd?
			ze := zettel.External{
				Sku: sku.External[kennung.Hinweis, *kennung.Hinweis]{
					FDs: sku.ExternalFDs{
						Objekte: kennung.FD{
							Path: z.Sku.Kennung.String(),
						},
					},
				},
			}

			if err1 := s.readZettelFromFile(&ze); err1 == nil {
				z.Objekte = ze.Objekte
				z.Sku.ObjekteSha = ze.Sku.ObjekteSha // TODO-P1 determine what else in sku is needed

				z.Verzeichnisse.ResetWithObjekte(z.Objekte)

				if err = w1(z); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			}

			if err = w1(z); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	return s.storeObjekten.Zettel().ReadAllSchwanzen(
		w,
	)
}

func (s *Store) ReadManyHistory(
	w schnittstellen.FuncIter[*zettel.Transacted],
) (err error) {
	queries := []func(schnittstellen.FuncIter[*zettel.Transacted]) error{
		s.storeObjekten.Zettel().ReadAll,
	}

	if s.erworben.IncludeCwd {
		queries = append(
			queries,
			func(w schnittstellen.FuncIter[*zettel.Transacted]) (err error) {
				var pz cwd.CwdFiles

				if pz, err = cwd.MakeCwdFilesAll(s.erworben, s.Standort.Cwd()); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = pz.Zettelen.Each(
					func(p cwd.Zettel) (err error) {
						var checked_out zettel.CheckedOut

						var readFunc func() (zettel.CheckedOut, error)

						switch {
						case p.GetAkteFD().Path == "":
							readFunc = func() (zettel.CheckedOut, error) {
								return s.Read(p.GetObjekteFD().Path)
							}

						case p.GetObjekteFD().Path == "":
							readFunc = func() (zettel.CheckedOut, error) {
								return s.ReadExternalZettelFromAktePath(p.GetAkteFD().Path)
							}

						default:
							// TODO-P3 validate that the zettel file points to the akte in the metadatei
							readFunc = func() (zettel.CheckedOut, error) {
								return s.Read(p.GetObjekteFD().Path)
							}
						}

						if checked_out, err = readFunc(); err != nil {
							// TODO-P3 decide if error handling like this is ok
							if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
								errors.Err().Printf("external zettel does not exist: %s", p)
							} else {
								errors.Err().Print(err)
							}

							err = nil
							return
						}

						zt := &zettel.Transacted{
							Sku:     checked_out.External.Sku.Transacted(),
							Objekte: checked_out.External.Objekte,
						}

						zt.Sku.Schwanz = s.sonnenaufgang
						zt.Verzeichnisse.ResetWithObjekte(zt.Objekte)

						if err = w(zt); err != nil {
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
			},
		)
	}

	return iter.Multiplex(
		w,
		queries...,
	)
}

func (s *Store) ReadFiles(
	fs cwd.CwdFiles,
	ms kennung.MetaSet,
	f schnittstellen.FuncIter[objekte.CheckedOutLike],
) (err error) {
	if err = s.storeObjekten.Query(
		ms,
		iter.MakeChain(
			func(e objekte.TransactedLike) (err error) {
				switch et := e.(type) {
				case *zettel.Transacted:
					var zco zettel.CheckedOut
					ok := false

					var ze cwd.Zettel

					if ze, ok = fs.GetZettel(et.Sku.Kennung); !ok {
						err = iter.MakeErrStopIteration()
						return
					}

					zco.External.Sku.FDs = ze.FDs
					zco.External.Sku.Kennung = ze.Kennung

					if err = s.readZettelFromFile(&zco.External); err != nil {
						if errors.IsNotExist(err) {
							err = iter.MakeErrStopIteration()
						} else {
							err = errors.Wrap(err)
						}

						return
					}

					zco.Internal = *et

					zco.DetermineState()

					if err = f(&zco); err != nil {
						err = errors.Wrap(err)
						return
					}

				case *typ.Transacted:
					var tco typ.CheckedOut
					ok := false

					var te cwd.Typ

					if te, ok = fs.GetTyp(et.Sku.Kennung); !ok {
						err = iter.MakeErrStopIteration()
						return
					}

					if tco.External, err = s.ReadTyp(te); err != nil {
						if errors.IsNotExist(err) {
							err = iter.MakeErrStopIteration()
						} else {
							err = errors.Wrapf(err, "CwdTyp: %#v", te)
						}

						return
					}

					tco.Internal = *et

					tco.DetermineState()

					if err = f(&tco); err != nil {
						err = errors.Wrap(err)
						return
					}

				case *etikett.Transacted:
					var tco etikett.CheckedOut
					ok := false

					var te cwd.Etikett

					if te, ok = fs.GetEtikett(et.Sku.Kennung); !ok {
						err = iter.MakeErrStopIteration()
						return
					}

					if tco.External, err = s.ReadEtikett(te); err != nil {
						if errors.IsNotExist(err) {
							err = iter.MakeErrStopIteration()
						} else {
							err = errors.Wrapf(err, "CwdEtikett: %#v", te)
						}

						return
					}

					tco.Internal = *et

					tco.DetermineState()

					if err = f(&tco); err != nil {
						err = errors.Wrap(err)
						return
					}

				default:
					err = errors.Implement()
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
				case cwd.Typ:
					todo.Implement()

				case cwd.Etikett:
					if err = s.storeObjekten.GetAbbrStore().EtikettExists(
						il.Kennung,
					); err == nil {
						err = iter.MakeErrStopIteration()
						return
					}

					err = nil

					var tco etikett.CheckedOut

					if tco.External, err = s.ReadEtikett(il); err != nil {
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

func (s *Store) readOneFS(p string) (cz zettel.CheckedOut, err error) {
	if cz.External, err = s.MakeExternalZettelFromZettel(p); err != nil {
		err = errors.Wrapf(err, "%s", p)
		return
	}

	if err = s.readZettelFromFile(&cz.External); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrapf(err, "%s", p)
			return
		}
	}

	return
}

func (s *Store) Read(p string) (cz zettel.CheckedOut, err error) {
	if cz, err = s.readOneFS(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	var zt *zettel.Transacted

	if zt, err = s.storeObjekten.Zettel().ReadOne(
		cz.External.Sku.Kennung,
	); err != nil {
		// if errors.Is(err, store_objekten.ErrNotFound{}) {
		// 	err = nil
		// } else {
		err = errors.Wrap(err)
		// }

		return
	}

	cz.Internal = *zt
	cz.DetermineState()

	if cz.State > objekte.CheckedOutStateExistsAndSame {
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
	}

	return
}
