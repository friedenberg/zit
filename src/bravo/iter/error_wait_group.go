package iter

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type ErrorWaitGroup interface {
	Do(schnittstellen.FuncError) bool
	DoAfter(schnittstellen.FuncError)
	Wait()
	GetError() error
}

func MakeErrorWaitGroup() ErrorWaitGroup {
	wg := &errorWaitGroup{
		lock:    &sync.RWMutex{},
		chDone:  make(chan struct{}),
		err:     errors.MakeMulti(),
		doAfter: make([]schnittstellen.FuncError, 0),
	}

	return wg
}

type errorWaitGroup struct {
	lock       *sync.RWMutex
	waitingFor int
	chDone     chan struct{}
	err        errors.Multi
	doAfter    []schnittstellen.FuncError
}

func (wg *errorWaitGroup) GetError() (err error) {
	wg.Wait()

	me := errors.MakeMulti(wg.err)

	defer func() {
		if !me.Empty() {
			err = me
		}
	}()

	for i := len(wg.doAfter) - 1; i >= 0; i-- {
		me.Add(wg.doAfter[i]())
	}

	return
}

func ErrorWaitGroupApply[T any](
	wg ErrorWaitGroup,
	s schnittstellen.Set[T],
	f schnittstellen.FuncIter[T],
) (d bool) {
	if err := s.Each(
		func(e T) (err error) {
			if wg.Do(
				func() error {
					return f(e)
				},
			) {
				err = MakeErrStopIteration()
			}

			return
		},
	); err != nil {
		d = true
	}

	return
}

func (wg *errorWaitGroup) Do(f schnittstellen.FuncError) (d bool) {
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

func (wg *errorWaitGroup) DoAfter(f schnittstellen.FuncError) {
	wg.lock.Lock()
	defer wg.lock.Unlock()

	wg.doAfter = append(wg.doAfter, f)

	return
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
