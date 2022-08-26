package store_working_directory

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/files"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/file_lock"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/golf/zettel_formats"
	"github.com/friedenberg/zit/src/hotel/collections"
)

type StoreZettel interface {
	Read(id id.Id) (z stored_zettel.Transacted, err error)
	ReadAkteSha(sha.Sha) (collections.SetTransacted, error)
	WriteZettelObjekte(z zettel.Zettel) (sh sha.Sha, err error)
	zettel.AkteWriterFactory
	zettel.AkteReaderFactory
}

type Store struct {
	lock *file_lock.Lock
	Konfig
	format       zettel.Format
	storeZettel  StoreZettel
	path         string
	cwd          string
	entries      map[string]Entry
	indexWasRead bool
	hasChanges   bool
}

func New(k Konfig, p string, storeZettel StoreZettel) (s *Store, err error) {
	s = &Store{
		Konfig:      k,
		format:      zettel_formats.Text{},
		storeZettel: storeZettel,
		path:        p,
		entries:     make(map[string]Entry),
	}

	s.lock = file_lock.New(path.Join(p, ".ZitCheckoutStoreLock"))

	if err = s.lock.Lock(); err != nil {
		err = errors.Error(err)
		return
	}

	if s.cwd, err = os.Getwd(); err != nil {
		err = errors.Error(err)
		return
	}

	logz.Print()

	return
}

func (s Store) IndexFilePath() string {
	return path.Join(s.path, ".ZitCheckoutStoreIndex")
}

func (s Store) flushToTemp() (tfp string, err error) {
	var f *os.File

	if f, err = open_file_guard.TempFile(); err != nil {
		err = errors.Error(err)
		return
	}

	tfp = f.Name()

	defer open_file_guard.Close(f)

	w := bufio.NewWriter(f)
	defer w.Flush()

	for p, e := range s.entries {
		out := fmt.Sprintf("%s %s\n", p, e)
		logz.Printf("flushing zettel: %q", out)
		w.WriteString(fmt.Sprint(out))
	}

	return
}

func (s Store) Flush() (err error) {
	if s.hasChanges {
		var tfp string

		if tfp, err = s.flushToTemp(); err != nil {
			err = errors.Error(err)
			return
		}

		logz.Printf("renaming %s to %s", tfp, s.IndexFilePath())
		if err = os.Rename(tfp, s.IndexFilePath()); err != nil {
			err = errors.Error(err)
			return
		}
	}

	if err = s.lock.Unlock(); err != nil {
		err = errors.Error(err)
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
		err = errors.Error(err)
		return
	}

	for _, p := range possible.Zettelen {
		if err = s.syncOne(p); err != nil {
			err = errors.Error(err)
			return
		}
	}

	s.indexWasRead = true

	return
}

func (s *Store) syncOne(p string) (err error) {
	logz.Output(2, fmt.Sprintln("will sync one: ", p))
	var hasCache, hasFs bool

	var fi os.FileInfo

	if fi, err = os.Stat(p); err != nil {
		if !os.IsNotExist(err) {
			err = errors.Error(err)
			return
		}
	} else {
		hasFs = true
	}

	var cached Entry

	cached, hasCache = s.entries[p]

	if !hasCache && !hasFs {
		logz.Print(p, ": no cache, no fs")
		return
	} else if hasCache {
		logz.Print(p, ": cache, no fs: deleting")
		delete(s.entries, p)
	} else {
		logz.Print(p, ": cache, fs")
		if !hasCache || fi.ModTime().After(cached.Time) {
			var ez stored_zettel.External

			if ez, err = s.MakeExternalZettelFromZettel(p); err != nil {
				err = errors.Error(err)
				return
			}

			if err = s.readZettelFromFile(&ez); err != nil {
				err = errors.Error(err)
				return
			}

			s.entries[p] = Entry{
				Time: fi.ModTime(),
				Sha:  ez.Stored.Sha,
			}

			s.hasChanges = true
		}
	}

	return
}

func (s Store) MakeExternalZettelFromZettel(p string) (ez stored_zettel.External, err error) {
	if p, err = filepath.Abs(p); err != nil {
		err = errors.Error(err)
		return
	}

	if p, err = filepath.Rel(s.path, p); err != nil {
		err = errors.Error(err)
		return
	}

	ez.Path = p

	head, tail := id.HeadTailFromFileName(p)

	if ez.Hinweis, err = hinweis.MakeBlindHinweis(head + "/" + tail); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s Store) readZettelFromFile(ez *stored_zettel.External) (err error) {
	logz.PrintDebug(ez)
	if !files.Exists(ez.Path) {
		//if the path does not have an extension, try looking for a file with that
		//extension
		//TODO modify this to use globs
		if filepath.Ext(ez.Path) == "" {
			ez.Path = ez.Path + ".md"
			return s.readZettelFromFile(ez)
		}

		err = os.ErrNotExist

		return
	}

	c := zettel.FormatContextRead{
		AkteWriterFactory: s.storeZettel,
	}

	var f *os.File

	if f, err = open_file_guard.Open(ez.Path); err != nil {
		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(f)

	c.In = f

	if _, err = s.format.ReadFrom(&c); err != nil {
		err = errors.Wrapped(err, "%s", f.Name())
		return
	}

	if ez.Stored.Sha, err = s.storeZettel.WriteZettelObjekte(c.Zettel); err != nil {
		err = errors.Wrapped(err, "%s", f.Name())
		return
	}

	ez.Named.Stored.Zettel = c.Zettel
	ez.AktePath = c.AktePath

	return
}

func (s *Store) Read(p string) (ez stored_zettel.External, err error) {
	if ez, err = s.MakeExternalZettelFromZettel(p); err != nil {
		err = errors.Wrapped(err, "%s", p)
		return
	}

	logz.Print(ez.Hinweis)

	if s.CacheEnabled {
		if err = s.ReadAll(); err != nil {
			err = errors.Wrapped(err, "%s", p)
			return
		}

		var cached Entry
		var hasEntry bool

		if cached, hasEntry = s.entries[p]; !hasEntry {
			logz.Printf("cached not found: %s\n", p)
			logz.Printf("%#v", s.entries)
			err = ErrNotInIndex(nil)
			return
		}

		var named stored_zettel.Transacted

		if named, err = s.storeZettel.Read(cached.Sha); err != nil {
			err = errors.Error(err)
			return
		}

		ez.Stored.Sha = named.Named.Stored.Sha
		ez.Stored.Zettel = named.Named.Stored.Zettel
	} else {
		if err = s.readZettelFromFile(&ez); err != nil {
			err = errors.Wrapped(err, "%s", p)
			return
		}
	}

	return
}
