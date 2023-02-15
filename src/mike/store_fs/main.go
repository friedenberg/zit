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
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/hinweisen"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/juliett/cwd"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/zettel_external"
	"github.com/friedenberg/zit/src/lima/store_objekten"
	"github.com/friedenberg/zit/src/lima/zettel_checked_out"
)

type Store struct {
	sonnenaufgang ts.Time
	erworben      konfig.Compiled
	standort.Standort

	format zettel.ObjekteFormat

	storeObjekten *store_objekten.Store

	zettelExternalLogPrinter schnittstellen.FuncIter[*zettel_external.Zettel]

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

func (s Store) MakeExternalZettelFromZettel(p string) (ez zettel_external.Zettel, err error) {
	if p, err = filepath.Abs(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if p, err = filepath.Rel(s.Cwd(), p); err != nil {
		err = errors.Wrap(err)
		return
	}

	ez.ZettelFD.Path = p

	head, tail := id.HeadTailFromFileName(p)

	if ez.Sku.Kennung, err = kennung.MakeHinweis(head + "/" + tail); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Store) readZettelFromFile(ez *zettel_external.Zettel) (err error) {
	if !files.Exists(ez.ZettelFD.Path) {
		// if the path does not have an extension, try looking for a file with that
		// extension
		// TODO-P4 modify this to use globs
		if filepath.Ext(ez.ZettelFD.Path) == "" {
			ez.ZettelFD.Path = ez.ZettelFD.Path + s.erworben.GetZettelFileExtension()
			return s.readZettelFromFile(ez)
		}

		err = os.ErrNotExist

		return
	}

	c := zettel.ObjekteParserContext{}

	var f *os.File

	if f, err = files.Open(ez.ZettelFD.Path); err != nil {
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
	ez.AkteFD.Path = c.AktePath

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

	var checked_out zettel_checked_out.Zettel

	var readFunc func() (zettel_checked_out.Zettel, error)

	var p cwd.CwdZettel

	for _, p1 := range pz.Zettelen {
		p = p1
		break
	}

	switch {
	case p.Akte.Path == "":
		readFunc = func() (zettel_checked_out.Zettel, error) {
			return s.Read(p.Zettel.Path)
		}

	case p.Zettel.Path == "":
		readFunc = func() (zettel_checked_out.Zettel, error) {
			return s.ReadExternalZettelFromAktePath(p.Akte.Path)
		}

	default:
		// TODO-P3 validate that the zettel file points to the akte in the metadatei
		readFunc = func() (zettel_checked_out.Zettel, error) {
			return s.Read(p.Zettel.Path)
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
			ze := zettel_external.Zettel{
				ZettelFD: kennung.FD{
					Path: z.Sku.Kennung.String(),
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

				for _, p := range pz.Zettelen {
					var checked_out zettel_checked_out.Zettel

					var readFunc func() (zettel_checked_out.Zettel, error)

					switch {
					case p.Akte.Path == "":
						readFunc = func() (zettel_checked_out.Zettel, error) {
							return s.Read(p.Zettel.Path)
						}

					case p.Zettel.Path == "":
						readFunc = func() (zettel_checked_out.Zettel, error) {
							return s.ReadExternalZettelFromAktePath(p.Akte.Path)
						}

					default:
						// TODO-P3 validate that the zettel file points to the akte in the metadatei
						readFunc = func() (zettel_checked_out.Zettel, error) {
							return s.Read(p.Zettel.Path)
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
						continue
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
				}

				return
			},
		)
	}

	return collections.Multiplex(
		w,
		queries...,
	)
}

func (s *Store) readOneFS(p string) (cz zettel_checked_out.Zettel, err error) {
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

func (s *Store) Read(p string) (cz zettel_checked_out.Zettel, err error) {
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

	if cz.State > zettel_checked_out.StateExistsAndSame {
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
