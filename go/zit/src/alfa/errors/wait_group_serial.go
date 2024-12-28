package errors

import "sync"

func MakeWaitGroupSerial() ErrorWaitGroup {
	wg := &errorWaitGroupSerial{
		do:      make([]Func, 0),
		doAfter: make([]FuncErrorWithStackInfo, 0),
	}

	return wg
}

type errorWaitGroupSerial struct {
	lock    sync.Mutex
	do      []Func
	doAfter []FuncErrorWithStackInfo
	isDone  bool
}

func (wg *errorWaitGroupSerial) GetError() (err error) {
	wg.lock.Lock()
	defer wg.lock.Unlock()

	wg.isDone = true

	me := MakeMulti()

	for _, f := range wg.do {
		if err = f(); err != nil {
			me.Add(Wrap(err))
			break
		}
	}

	for i := len(wg.doAfter) - 1; i >= 0; i-- {
		doAfter := wg.doAfter[i]
		err := doAfter.Func()
		if err != nil {
			me.Add(doAfter.Wrap(err))
		}
	}

	err = me.GetError()

	return
}

func (wg *errorWaitGroupSerial) Do(f Func) (d bool) {
	wg.lock.Lock()
	defer wg.lock.Unlock()

	if wg.isDone {
		return false
	}

	wg.do = append(wg.do, f)

	return true
}

func (wg *errorWaitGroupSerial) DoAfter(f Func) {
	wg.lock.Lock()
	defer wg.lock.Unlock()

	wg.do = append(wg.do, f)

	return
}
