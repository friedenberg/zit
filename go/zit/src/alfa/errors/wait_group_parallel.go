package errors

import "sync"

func MakeWaitGroupParallel() WaitGroup {
	wg := &waitGroupParallel{
		lock:    &sync.Mutex{},
		inner:   &sync.WaitGroup{},
		err:     MakeMulti(),
		doAfter: make([]FuncWithStackInfo, 0),
	}

	return wg
}

type waitGroupParallel struct {
	lock    *sync.Mutex
	inner   *sync.WaitGroup
	err     Multi
	doAfter []FuncWithStackInfo

	isDone bool
}

func (wg *waitGroupParallel) GetError() (err error) {
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

func (wg *waitGroupParallel) Do(f Func) (d bool) {
	wg.lock.Lock()

	if wg.isDone {
		wg.lock.Unlock()
		return false
	}

	wg.lock.Unlock()

	wg.inner.Add(1)

	si, _ := MakeStackFrame(1)

	go func() {
		err := f()
		wg.doneWith(si, err)
	}()

	return true
}

func (wg *waitGroupParallel) DoAfter(f Func) {
	wg.lock.Lock()
	defer wg.lock.Unlock()

	si, _ := MakeStackFrame(1)

	wg.doAfter = append(
		wg.doAfter,
		FuncWithStackInfo{
			Func:      f,
			StackFrame: si,
		},
	)
}

func (wg *waitGroupParallel) doneWith(si StackFrame, err error) {
	wg.inner.Done()

	if err != nil {
		wg.err.Add(si.Wrap(err))
	}
}

func (wg *waitGroupParallel) wait() {
	wg.inner.Wait()

	wg.lock.Lock()
	defer wg.lock.Unlock()

	wg.isDone = true
}
