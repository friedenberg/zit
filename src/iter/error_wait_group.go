package iter

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type ErrorWaitGroup interface {
	Do(func() error) bool
	Wait()
	GetError() error
}

func MakeErrorWaitGroup() ErrorWaitGroup {
	wg := &errorWaitGroup{
		lock:   &sync.RWMutex{},
		chDone: make(chan struct{}),
		err:    errors.MakeMulti(),
	}

	return wg
}

type errorWaitGroup struct {
	lock       *sync.RWMutex
	waitingFor int
	chDone     chan struct{}
	err        errors.Multi
}

func (wg *errorWaitGroup) GetError() (err error) {
	wg.Wait()

	if !wg.err.Empty() {
		err = wg.err
	}

	return
}

func (wg *errorWaitGroup) Do(f func() error) (d bool) {
	wg.lock.Lock()
	defer wg.lock.Unlock()

	if wg.isDone() {
		return true
	}

	wg.waitingFor += 1

	go func() {
		wg.doneWith(f())
	}()

	return false
}

func (wg *errorWaitGroup) doneWith(err error) {
	wg.lock.Lock()
	defer wg.lock.Unlock()

	if err != nil {
		wg.err.Add(err)
		wg.waitingFor = 0
	} else {
		wg.waitingFor--
	}

	if wg.waitingFor == 0 {
		close(wg.chDone)
	}
}

func (wg *errorWaitGroup) isDone() bool {
	select {
	case <-wg.chDone:
		return true

	default:
		return false
	}
}

func (wg *errorWaitGroup) IsDone() bool {
	wg.lock.RLock()
	defer wg.lock.RUnlock()

	return wg.isDone()
}

func (wg *errorWaitGroup) Wait() {
	<-wg.chDone
}
