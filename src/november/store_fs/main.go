package store_fs

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/standort"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/hinweisen"
	"github.com/friedenberg/zit/src/foxtrot/id"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/fd"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/juliett/konfig_compiled"
	"github.com/friedenberg/zit/src/kilo/cwd_files"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/store_objekten"
	"github.com/friedenberg/zit/src/mike/zettel_checked_out"
)

type Store struct {
	sonnenaufgang ts.Time
	konfig        konfig_compiled.Compiled
	standort.Standort

	format zettel.Format

	storeObjekten *store_objekten.Store

	zettelCheckedOutWriters ZettelCheckedOutLogWriters

	entries      map[string]Entry
	indexWasRead bool
	hasChanges   bool
}

func New(
	t ts.Time,
	k konfig_compiled.Compiled,
	st standort.Standort,
	storeObjekten *store_objekten.Store,
) (s *Store, err error) {
	s = &Store{
		sonnenaufgang: t,
		konfig:        k,
		Standort:      st,
		format: zettel.MakeTextParser(
			storeObjekten,
			nil,
		),
		storeObjekten: storeObjekten,
		entries:       make(map[string]Entry),
	}

	return
}

func (s *Store) SetZettelCheckedOutWriters(
	zcow ZettelCheckedOutLogWriters,
) {
	s.zettelCheckedOutWriters = zcow
}

// TODO-P3 move to standort
func (s Store) IndexFilePath() string {
	return path.Join(s.Cwd(), ".ZitCheckoutStoreIndex")
}

func (s Store) flushToTemp() (tfp string, err error) {
	var f *os.File

	if f, err = files.TempFile(); err != nil {
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

	if ez.Sku.Kennung, err = hinweis.Make(head + "/" + tail); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Store) readZettelFromFile(ez *zettel_external.Zettel) (err error) {
	if !files.Exists(ez.ZettelFD.Path) {
		//if the path does not have an extension, try looking for a file with that
		//extension
		//TODO-P4 modify this to use globs
		if filepath.Ext(ez.ZettelFD.Path) == "" {
			ez.ZettelFD.Path = ez.ZettelFD.Path + s.konfig.GetZettelFileExtension()
			return s.readZettelFromFile(ez)
		}

		err = os.ErrNotExist

		return
	}

	c := zettel.FormatContextRead{}

	var f *os.File

	if f, err = files.Open(ez.ZettelFD.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	c.In = f

	if _, err = s.format.ReadFrom(&c); err != nil {
		err = errors.Wrapf(err, "%s", f.Name())
		return
	}

	if ez.Sku.Sha, err = s.storeObjekten.Zettel().WriteZettelObjekte(
		c.Zettel,
	); err != nil {
		err = errors.Wrapf(err, "%s", f.Name())
		return
	}

	ez.Objekte = c.Zettel
	ez.AkteFD.Path = c.AktePath

	var unrecoverableErrors errors.Multi

	for _, e := range c.RecoverableErrors {
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

			if mutter, err = s.storeObjekten.Zettel().ReadHinweisSchwanzen(ez.Sku.Kennung); err != nil {
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
// TODO-P1 transition to ZettelVerzeichnisse
func (s *Store) ReadOne(h hinweis.Hinweis) (zt *zettel.Transacted, err error) {
	if zt, err = s.storeObjekten.Zettel().ReadHinweisSchwanzen(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !s.konfig.IncludeCwd {
		return
	}

	var pz cwd_files.CwdFiles

	if pz, err = cwd_files.MakeCwdFilesExactly(
		s.konfig,
		s.Standort.Cwd(),
		fmt.Sprintf("%s.%s", h, s.konfig.FileExtensions.Zettel),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var checked_out zettel_checked_out.Zettel

	var readFunc func() (zettel_checked_out.Zettel, error)

	var p cwd_files.CwdZettel

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
		//TODO-P3 validate that the zettel file points to the akte in the metadatei
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
	w1 collections.WriterFunc[*zettel.Verzeichnisse],
) (err error) {
	w := w1

	if s.konfig.IncludeCwd {
		w = func(z *zettel.Verzeichnisse) (err error) {
			//TODO-P2 akte fd?
			ze := zettel_external.Zettel{
				ZettelFD: fd.FD{
					Path: z.Transacted.Sku.Kennung.String(),
				},
			}

			if err1 := s.readZettelFromFile(&ze); err1 == nil {
				z1 := &zettel.Transacted{
					Sku:     z.Transacted.Sku,
					Objekte: ze.Objekte,
				}

				z1.Sku.Sha = ze.Sku.Sha

				z2 := &zettel.Verzeichnisse{}

				z2.ResetWithTransacted(z1)

				if err = w1(z2); err != nil {
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

	return s.storeObjekten.Zettel().ReadAllSchwanzenVerzeichnisse(
		w,
	)
}

func (s *Store) ReadManyHistory(
	w collections.WriterFunc[*zettel.Verzeichnisse],
) (err error) {
	queries := []func(collections.WriterFunc[*zettel.Verzeichnisse]) error{
		s.storeObjekten.Zettel().ReadAllVerzeichnisse,
	}

	if s.konfig.IncludeCwd {
		queries = append(
			queries,
			func(w collections.WriterFunc[*zettel.Verzeichnisse]) (err error) {
				var pz cwd_files.CwdFiles

				if pz, err = cwd_files.MakeCwdFilesAll(s.konfig, s.Standort.Cwd()); err != nil {
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
						//TODO-P3 validate that the zettel file points to the akte in the metadatei
						readFunc = func() (zettel_checked_out.Zettel, error) {
							return s.Read(p.Zettel.Path)
						}
					}

					if checked_out, err = readFunc(); err != nil {
						//TODO-P3 decide if error handling like this is ok
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

					zv := &zettel.Verzeichnisse{}
					zv.ResetWithTransacted(zt)

					if err = w(zv); err != nil {
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

// TODO-P1 deprecate in favor of above method
func (s *Store) ZettelTransactedWriter(
	w1 collections.WriterFunc[*zettel.Transacted],
) (w collections.WriterFunc[*zettel.Transacted]) {
	return func(z *zettel.Transacted) (err error) {
		//TODO-P2 akte fd?
		ze := zettel_external.Zettel{
			ZettelFD: fd.FD{
				Path: z.Sku.Kennung.String(),
			},
		}

		if err1 := s.readZettelFromFile(&ze); err1 == nil {
			z1 := &zettel.Transacted{
				Sku:     z.Sku,
				Objekte: ze.Objekte,
			}

			z1.Sku.Sha = ze.Sku.Sha

			if err = w1(z1); err != nil {
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

func (s *Store) Read(p string) (cz zettel_checked_out.Zettel, err error) {
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

	var zt *zettel.Transacted

	if zt, err = s.storeObjekten.Zettel().ReadHinweisSchwanzen(
		cz.External.Sku.Kennung,
	); err != nil {
		if errors.Is(err, store_objekten.ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	cz.Internal = *zt
	cz.DetermineState()

	if cz.State > zettel_checked_out.StateExistsAndSame {
		//TODO-P4 rewrite with verzeichnisseAll
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
