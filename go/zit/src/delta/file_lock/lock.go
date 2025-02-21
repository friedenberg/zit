package file_lock

import (
	"io/fs"
	"os"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
)

type Lock struct {
	envUI env_ui.Env
	path  string
	mutex sync.Mutex
	file  *os.File
}

func New(envUI env_ui.Env, path string) (l *Lock) {
	return &Lock{
		envUI: envUI,
		path:  path,
	}
}

func (l *Lock) Path() string {
	return l.path
}

func (l *Lock) IsAcquired() (acquired bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	acquired = l.file != nil

	return
}

func (lock *Lock) Lock() (err error) {
	if !lock.mutex.TryLock() {
		err = errors.Errorf("attempting concurrent locks")
		return
	}

	defer lock.mutex.Unlock()

	if lock.file != nil {
		err = errors.Errorf("already locked")
		return
	}

	createLock := func(path string) (*os.File, error) {
		return files.OpenFile(
			path,
			os.O_RDONLY|os.O_EXCL|os.O_CREATE,
			0o755,
		)
	}

	if lock.file, err = files.MakeDirIfNecessary(
		lock.Path(),
		createLock,
	); err != nil {
		if errors.Is(err, fs.ErrExist) {
			err = ErrUnableToAcquireLock{envUI: lock.envUI, Path: lock.Path()}
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	return
}

func (lock *Lock) Unlock() (err error) {
	if !lock.mutex.TryLock() {
		err = errors.Errorf("attempting concurrent locks")
		return
	}

	defer lock.mutex.Unlock()

	if err = files.Close(lock.file); err != nil {
		err = errors.Wrapf(err, "File: %v", lock.file)
		return
	}

	lock.file = nil

	if err = os.Remove(lock.Path()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
