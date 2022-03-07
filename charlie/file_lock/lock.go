package file_lock

import (
	"os"
	"sync"
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

	if l.f, err = _OpenFile(l.Path(), os.O_RDONLY|os.O_EXCL|os.O_CREATE, 755); err != nil {
		err = _Errorf("lockfile already exists, unable to acquire lock: %w", err)
		return
	}

	return
}

func (l *Lock) Unlock() (err error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if err = _Close(l.f); err != nil {
		err = _Error(err)
		return
	}

	l.f = nil

	if err = os.Remove(l.Path()); err != nil {
		err = _Error(err)
		return
	}

	return
}
