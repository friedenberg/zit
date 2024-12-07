package file_lock

import (
	"io/fs"
	"os"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
)

type Lock struct {
	path  string
	mutex sync.Locker
	f     *os.File
}

func New(path string) (l *Lock) {
	return &Lock{
		path:  path,
		mutex: &sync.Mutex{},
	}
}

func (l Lock) Path() string {
	return l.path
}

func (l *Lock) IsAcquired() (acquired bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	acquired = l.f != nil

	return
}

func (l *Lock) Lock() (err error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.f != nil {
		err = errors.Errorf("already locked")
		return
	}

	ui.Log().Caller(2, "locking "+l.Path())
	if l.f, err = files.OpenFile(l.Path(), os.O_RDONLY|os.O_EXCL|os.O_CREATE, 755); err != nil {
		if errors.Is(err, fs.ErrExist) {
			err = errors.Wrapf(err, "lockfile already exists, unable to acquire lock: %s", l.Path())
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	return
}

func (l *Lock) Unlock() (err error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// TODO-P4 determine if there's some way for error.Deferred to correctly log
	// the location of this
	ui.Log().Caller(2, "unlocking "+l.Path())
	if err = files.Close(l.f); err != nil {
		err = errors.Wrapf(err, "File: %v", l.f)
		return
	}

	l.f = nil

	if err = os.Remove(l.Path()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
