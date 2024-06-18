package iter

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
)

type ErrorWaitGroup interface {
	Do(schnittstellen.FuncError) bool
	DoAfter(schnittstellen.FuncError)
	GetError() error
}

func MakeErrorWaitGroupParallel() ErrorWaitGroup {
	wg := &errorWaitGroupParallel{
		lock:    &sync.Mutex{},
		inner:   &sync.WaitGroup{},
		err:     errors.MakeMulti(),
		doAfter: make([]schnittstellen.FuncError, 0),
	}

	return wg
}

type errorWaitGroupParallel struct {
	lock    *sync.Mutex
	inner   *sync.WaitGroup
	err     errors.Multi
	doAfter []schnittstellen.FuncError

	isDone bool
}

func (wg *errorWaitGroupParallel) GetError() (err error) {
	wg.wait()

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
	s schnittstellen.SetLike[T],
	f schnittstellen.FuncIter[T],
) (d bool) {
	if err := s.Each(
		func(e T) (err error) {
			if !wg.Do(
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

func (wg *errorWaitGroupParallel) Do(f schnittstellen.FuncError) (d bool) {
	wg.lock.Lock()

	if wg.isDone {
		wg.lock.Unlock()
		return false
	}

	wg.lock.Unlock()

	wg.inner.Add(1)

	go func() {
		wg.doneWith(f())
	}()

	return true
}

func (wg *errorWaitGroupParallel) DoAfter(f schnittstellen.FuncError) {
	wg.lock.Lock()
	defer wg.lock.Unlock()

	wg.doAfter = append(wg.doAfter, f)

	return
}

func (wg *errorWaitGroupParallel) doneWith(err error) {
	wg.inner.Done()

	if err != nil {
		wg.err.Add(err)
	}
}

func (wg *errorWaitGroupParallel) wait() {
	wg.inner.Wait()

	wg.lock.Lock()
	defer wg.lock.Unlock()

	wg.isDone = true
}
