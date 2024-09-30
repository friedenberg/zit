package quiter

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func MakeErrorWaitGroupSerial() ErrorWaitGroup {
	wg := &errorWaitGroupSerial{
		do:      make([]interfaces.FuncError, 0),
		doAfter: make([]interfaces.FuncError, 0),
	}

	return wg
}

type errorWaitGroupSerial struct {
	lock        sync.Mutex
	do, doAfter []interfaces.FuncError
	isDone      bool
}

func (wg *errorWaitGroupSerial) GetError() (err error) {
	wg.lock.Lock()
	defer wg.lock.Unlock()

	wg.isDone = true

	me := errors.MakeMulti()

	for _, f := range wg.do {
		if err = f(); err != nil {
			me.Add(errors.Wrap(err))
			break
		}
	}

	for _, f := range wg.doAfter {
		if err = f(); err != nil {
			me.Add(errors.Wrap(err))
		}
	}

	err = me.GetError()

	return
}

func (wg *errorWaitGroupSerial) Do(f interfaces.FuncError) (d bool) {
	wg.lock.Lock()
	defer wg.lock.Unlock()

	if wg.isDone {
		return false
	}

	wg.do = append(wg.do, f)

	return true
}

func (wg *errorWaitGroupSerial) DoAfter(f interfaces.FuncError) {
	wg.lock.Lock()
	defer wg.lock.Unlock()

	wg.do = append(wg.do, f)

	return
}
