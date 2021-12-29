package file_lock

import (
	"os"
	"sync"
)

type Lock struct {
	path     string
	mutex    sync.Locker
	acquired bool
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

func (l Lock) IsAcquired() (acquired bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	acquired = l.acquired

	return
}

func (l Lock) Lock() (err error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	var f *os.File

	if f, err = _OpenFile(l.Path(), os.O_RDONLY|os.O_EXCL|os.O_CREATE, 755); err != nil {
		err = _Errorf("lockfile already exists, unable to acquire lock: %w", err)
		return
	}

	l.acquired = true

	err = _Close(f)

	return
}

func (l Lock) Unlock() (err error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	err = os.Remove(l.Path())
	l.acquired = false

	return
}
