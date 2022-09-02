package store_working_directory

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/open_file_guard"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/file_lock"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/id"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/store_objekten"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
)

type StoreZettel interface {
	Read(id id.Id) (z zettel_transacted.Zettel, err error)
	ReadAkteSha(sha.Sha) (zettel_transacted.Set, error)
	ReadZettelSha(sha.Sha) (zettel_transacted.Set, error)
	ReadBezeichnung(string) (zettel_transacted.Set, error)
	WriteZettelObjekte(z zettel.Zettel) (sh sha.Sha, err error)
	zettel.AkteWriterFactory
	zettel.AkteReaderFactory
}

type Store struct {
	lock *file_lock.Lock
	Konfig
	format        zettel.Format
	storeObjekten StoreZettel
	path          string
	cwd           string
	entries       map[string]Entry
	indexWasRead  bool
	hasChanges    bool
}

func New(k Konfig, p string, storeObjekten StoreZettel) (s *Store, err error) {
	s = &Store{
		Konfig:        k,
		format:        zettel.Text{},
		storeObjekten: storeObjekten,
		path:          p,
		entries:       make(map[string]Entry),
	}

	s.lock = file_lock.New(path.Join(p, ".ZitCheckoutStoreLock"))

	if err = s.lock.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.cwd, err = os.Getwd(); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Print()

	return
}

func (s Store) IndexFilePath() string {
	return path.Join(s.path, ".ZitCheckoutStoreIndex")
}

func (s Store) flushToTemp() (tfp string, err error) {
	var f *os.File

	if f, err = open_file_guard.TempFile(); err != nil {
		err = errors.Wrap(err)
		return
	}

	tfp = f.Name()

	defer open_file_guard.Close(f)

	w := bufio.NewWriter(f)
	defer w.Flush()

	for p, e := range s.entries {
		out := fmt.Sprintf("%s %s\n", p, e)
		errors.Printf("flushing zettel: %q", out)
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

		errors.Printf("renaming %s to %s", tfp, s.IndexFilePath())
		if err = os.Rename(tfp, s.IndexFilePath()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.lock.Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadAll() (err error) {
	if s.indexWasRead {
		return
	}

	var possible CwdFiles

	if possible, err = s.GetPossibleZettels(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, p := range possible.Zettelen {
		if err = s.syncOne(p); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	s.indexWasRead = true

	return
}

func (s *Store) syncOne(p string) (err error) {
	errors.Output(2, fmt.Sprintln("will sync one: ", p))
	var hasCache, hasFs bool

	var fi os.FileInfo

	if fi, err = os.Stat(p); err != nil {
		if !os.IsNotExist(err) {
			err = errors.Wrap(err)
			return
		}
	} else {
		hasFs = true
	}

	var cached Entry

	cached, hasCache = s.entries[p]

	if !hasCache && !hasFs {
		errors.Print(p, ": no cache, no fs")
		return
	} else if hasCache {
		errors.Print(p, ": cache, no fs: deleting")
		delete(s.entries, p)
	} else {
		errors.Print(p, ": cache, fs")
		if !hasCache || fi.ModTime().After(cached.Time) {
			var ez zettel_external.Zettel

			if ez, err = s.MakeExternalZettelFromZettel(p); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = s.readZettelFromFile(&ez); err != nil {
				err = errors.Wrap(err)
				return
			}

			s.entries[p] = Entry{
				Time: fi.ModTime(),
				Sha:  ez.Named.Stored.Sha,
			}

			s.hasChanges = true
		}
	}

	return
}

func (s Store) MakeExternalZettelFromZettel(p string) (ez zettel_external.Zettel, err error) {
	if p, err = filepath.Abs(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if p, err = filepath.Rel(s.path, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	ez.ZettelFD.Path = p

	head, tail := id.HeadTailFromFileName(p)

	if ez.Named.Hinweis, err = hinweis.Make(head + "/" + tail); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Store) readZettelFromFile(ez *zettel_external.Zettel) (err error) {
	errors.PrintDebug(ez)
	if !files.Exists(ez.ZettelFD.Path) {
		//if the path does not have an extension, try looking for a file with that
		//extension
		//TODO modify this to use globs
		if filepath.Ext(ez.ZettelFD.Path) == "" {
			ez.ZettelFD.Path = ez.ZettelFD.Path + ".md"
			return s.readZettelFromFile(ez)
		}

		err = os.ErrNotExist

		return
	}

	c := zettel.FormatContextRead{
		AkteWriterFactory: s.storeObjekten,
	}

	var f *os.File

	if f, err = open_file_guard.Open(ez.ZettelFD.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer open_file_guard.Close(f)

	c.In = f

	if _, err = s.format.ReadFrom(&c); err != nil {
		err = errors.Wrapf(err, "%s", f.Name())
		return
	}

	if ez.Named.Stored.Sha, err = s.storeObjekten.WriteZettelObjekte(c.Zettel); err != nil {
		err = errors.Wrapf(err, "%s", f.Name())
		return
	}

	ez.Named.Stored.Zettel = c.Zettel
	ez.AkteFD.Path = c.AktePath

	return
}

func (s *Store) Read(p string) (cz zettel_checked_out.Zettel, err error) {
	if cz.External, err = s.MakeExternalZettelFromZettel(p); err != nil {
		err = errors.Wrapf(err, "%s", p)
		return
	}

	if s.CacheEnabled {
		if err = s.ReadAll(); err != nil {
			err = errors.Wrapf(err, "%s", p)
			return
		}

		var cached Entry
		var hasEntry bool

		if cached, hasEntry = s.entries[p]; !hasEntry {
			errors.Printf("cached not found: %s", p)
			errors.Printf("%#v", s.entries)
			err = ErrNotInIndex(nil)
			return
		}

		var named zettel_transacted.Zettel

		if named, err = s.storeObjekten.Read(cached.Sha); err != nil {
			err = errors.Wrap(err)
			return
		}

		cz.External.Named.Stored.Sha = named.Named.Stored.Sha
		cz.External.Named.Stored.Zettel = named.Named.Stored.Zettel
	} else {
		if err = s.readZettelFromFile(&cz.External); err != nil {
			if errors.IsNotExist(err) {
				err = nil
			} else {
				err = errors.Wrapf(err, "%s", p)
				return
			}
		}

		if cz.Internal, err = s.storeObjekten.Read(cz.External.Named.Hinweis); err != nil {
			if errors.Is(err, store_objekten.ErrNotFound{}) {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		cz.DetermineState()

		if cz.State > zettel_checked_out.StateExistsAndSame {
			exSha := cz.External.Named.Stored.Sha
			cz.Matches.Zettelen, _ = s.storeObjekten.ReadZettelSha(exSha)

			exAkteSha := cz.External.Named.Stored.Zettel.Akte
			cz.Matches.Akten, _ = s.storeObjekten.ReadAkteSha(exAkteSha)

			bez := cz.External.Named.Stored.Zettel.Bezeichnung.String()
			cz.Matches.Bezeichnungen, _ = s.storeObjekten.ReadBezeichnung(bez)
		}

		return
	}

	return
}
