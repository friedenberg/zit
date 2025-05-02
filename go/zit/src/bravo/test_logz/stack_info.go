package test_logz

import (
	"runtime"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

type (
	StackInfo = errors.StackFrame
)

func MakeStackInfo(t *T, skip int) (si StackInfo) {
	var pc uintptr
	ok := false
	pc, _, _, ok = runtime.Caller(skip + 1)

	if !ok {
		t.Fatal("failed to make stack info")
	}

	frames := runtime.CallersFrames([]uintptr{pc})

	frame, _ := frames.Next()
	si = errors.MakeStackFrameFromFrame(frame)

	return
}
