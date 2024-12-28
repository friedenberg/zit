package errors

import "sync"

func MakeWaitGroupParallel() ErrorWaitGroup {
	wg := &errorWaitGroupParallel{
		lock:    &sync.Mutex{},
		inner:   &sync.WaitGroup{},
		err:     MakeMulti(),
		doAfter: make([]FuncErrorWithStackInfo, 0),
	}

	return wg
}

type FuncErrorWithStackInfo struct {
	Func
	StackInfo
}

type errorWaitGroupParallel struct {
	lock    *sync.Mutex
	inner   *sync.WaitGroup
	err     Multi
	doAfter []FuncErrorWithStackInfo

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
		doAfter := wg.doAfter[i]
		err := doAfter.Func()
		if err != nil {
			wg.err.Add(doAfter.Wrap(err))
		}
	}

	return
}

func (wg *errorWaitGroupParallel) Do(f Func) (d bool) {
	wg.lock.Lock()

	if wg.isDone {
		wg.lock.Unlock()
		return false
	}

	wg.lock.Unlock()

	wg.inner.Add(1)

	si, _ := MakeStackInfo(1)

	go func() {
		err := f()
		wg.doneWith(si, err)
	}()

	return true
}

func (wg *errorWaitGroupParallel) DoAfter(f Func) {
	wg.lock.Lock()
	defer wg.lock.Unlock()

	si, _ := MakeStackInfo(1)

	wg.doAfter = append(
		wg.doAfter,
		FuncErrorWithStackInfo{
			Func:      f,
			StackInfo: si,
		},
	)
}

func (wg *errorWaitGroupParallel) doneWith(si StackInfo, err error) {
	wg.inner.Done()

	if err != nil {
		wg.err.Add(si.Wrap(err))
	}
}

func (wg *errorWaitGroupParallel) wait() {
	wg.inner.Wait()

	wg.lock.Lock()
	defer wg.lock.Unlock()

	wg.isDone = true
}
