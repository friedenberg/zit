package checkout_store

import (
	"os"
	"path"

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
}

func New(p string, akteWriterFactory zettel.AkteWriterFactory) (s Store, err error) {
	s = Store{
		format:            zettel_formats.Text{},
		akteWriterFactory: akteWriterFactory,
		path:              p,
		entries:           make(map[string]Entry),
	}

	s.lock = file_lock.New(path.Join(p, ".CheckoutStoreLock"))

	if err = s.lock.Lock(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s Store) Flush() (err error) {
	//TODO flush index
	if err = s.lock.Unlock(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s Store) ReadAll() (err error) {
	return
}

func (s Store) readFromIndex(p string) (ez stored_zettel.External, err error) {
	// var fi os.FileInfo

	// if fi, err = os.Stat(p); err != nil {
	// 	if os.IsNotExist(err) {
	// 		err = ErrNotInIndex(err)
	// 	} else {
	// 		err = errors.Error(err)
	// 	}

	// 	return
	// }

	return
}

func (s Store) Read(p string) (ez stored_zettel.External, err error) {
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
