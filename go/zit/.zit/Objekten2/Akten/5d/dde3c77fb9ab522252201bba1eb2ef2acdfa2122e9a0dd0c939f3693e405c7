package quiter

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type ErrorWaitGroup interface {
	Do(interfaces.FuncError) bool
	DoAfter(interfaces.FuncError)
	GetError() error
}

func MakeErrorWaitGroupParallel() ErrorWaitGroup {
	wg := &errorWaitGroupParallel{
		lock:    &sync.Mutex{},
		inner:   &sync.WaitGroup{},
		err:     errors.MakeMulti(),
		doAfter: make([]interfaces.FuncError, 0),
	}

	return wg
}

type errorWaitGroupParallel struct {
	lock    *sync.Mutex
	inner   *sync.WaitGroup
	err     errors.Multi
	doAfter []interfaces.FuncError

	isDone bool
}

func (wg *errorWaitGroupParallel) GetError() (err error) {
	wg.wait()

	defer func() {
		if !wg.err.Empty() {
			err = wg.err
		}
	}()

	for i := len(wg.doAfter) - 1; i >= 0; i-- {
		wg.err.Add(wg.doAfter[i]())
	}

	return
}

func ErrorWaitGroupApply[T any](
	wg ErrorWaitGroup,
	s interfaces.SetLike[T],
	f interfaces.FuncIter[T],
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

func (wg *errorWaitGroupParallel) Do(f interfaces.FuncError) (d bool) {
	wg.lock.Lock()

	if wg.isDone {
		wg.lock.Unlock()
		return false
	}

	wg.lock.Unlock()

	wg.inner.Add(1)

	si, _ := errors.MakeStackInfo(2)

	go func() {
		err := f()
		wg.doneWith(si, err)
	}()

	return true
}

func (wg *errorWaitGroupParallel) DoAfter(f interfaces.FuncError) {
	wg.lock.Lock()
	defer wg.lock.Unlock()

	wg.doAfter = append(wg.doAfter, f)
}

func (wg *errorWaitGroupParallel) doneWith(si errors.StackInfo, err error) {
	wg.inner.Done()

	if err != nil {
		wg.err.Add(err)
		// wg.err.Add(si.Wrap(err))
	}
}

func (wg *errorWaitGroupParallel) wait() {
	wg.inner.Wait()

	wg.lock.Lock()
	defer wg.lock.Unlock()

	wg.isDone = true
}
