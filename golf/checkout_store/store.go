package checkout_store

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/files"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/file_lock"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
)

type Store struct {
	lock              *file_lock.Lock
	format            zettel.Format
	akteWriterFactory zettel.AkteWriterFactory
	path              string
	entries           map[string]Entry
	indexWasRead      bool
	hasChanges        bool
}

func New(p string, akteWriterFactory zettel.AkteWriterFactory) (s *Store, err error) {
	s = &Store{
		format:            zettel_formats.Text{},
		akteWriterFactory: akteWriterFactory,
		path:              p,
		entries:           make(map[string]Entry),
	}

	s.lock = file_lock.New(path.Join(p, ".ZitCheckoutStoreLock"))

	if err = s.lock.Lock(); err != nil {
		err = errors.Error(err)
		return
	}

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
		log.Printf("flushing zettel: %q", out)
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

		log.Printf("renaming %s to %s", tfp, s.IndexFilePath())
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

	var possible []string

	if possible, err = s.GetPossibleZettels(); err != nil {
		err = errors.Error(err)
		return
	}

	for _, p := range possible {
		if err = s.syncOne(p); err != nil {
			err = errors.Error(err)
			return
		}
	}

	s.indexWasRead = true

	return
}

func (s *Store) syncOne(p string) (err error) {
	log.Output(2, fmt.Sprintln("will sync one: ", p))
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
		log.Print(p, ": no cache, no fs")
		return
	} else if hasCache {
		log.Print(p, ": cache, no fs: deleting")
		delete(s.entries, p)
	} else {
		log.Print(p, ": cache, fs")
		if !hasCache ||
			fi.ModTime().After(cached.ZettelTime) ||
			fi.ModTime().After(cached.AkteTime) {
			var ez stored_zettel.External

			if ez, err = s.readZettelFromFile(p); err != nil {
				err = errors.Error(err)
				return
			}

			s.entries[p] = Entry{
				ZettelTime: fi.ModTime(),
				External:   ez,
			}

			s.hasChanges = true
		}
	}

	return
}

func (s Store) readZettelFromFile(p string) (ez stored_zettel.External, err error) {
	log.Print(p, ": reading from fs")
	ez.Path = p

	head, tail := id.HeadTailFromFileName(p)

	if ez.Hinweis, err = hinweis.MakeBlindHinweis(head + "/" + tail); err != nil {
		err = errors.Error(err)
		return
	}

	c := zettel.FormatContextRead{
		AkteWriterFactory: s.akteWriterFactory,
	}

	var f *os.File

	if !files.Exists(p) {
		err = os.ErrNotExist
		return
	}

	if f, err = os.Open(p); err != nil {
		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(f)

	c.In = f

	if _, err = s.format.ReadFrom(&c); err != nil {
		err = errors.Errorf("%s: %w", f.Name(), err)
		return
	}

	ez.Zettel = c.Zettel
	ez.AktePath = c.AktePath

	return
}

func (s *Store) Read(p string) (ez stored_zettel.External, err error) {
	if p, err = filepath.Abs(p); err != nil {
		err = errors.Error(err)
		return
	}

	if p, err = filepath.Rel(s.path, p); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.ReadAll(); err != nil {
		err = errors.Errorf("%w: %s", err, p)
		return
	}

	var cached Entry
	var hasEntry bool

	if cached, hasEntry = s.entries[p]; !hasEntry {
		log.Printf("cached not found: %s\n", p)
		log.Printf("%#v", s.entries)
		err = ErrNotInIndex(nil)
		return
	}

	ez = cached.External

	log.Print(ez)

	return
}
