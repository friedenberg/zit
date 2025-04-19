package debug

import (
	"os"
	"runtime/debug"
	"runtime/pprof"
	"runtime/trace"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
)

type Context struct {
	filePprofCpu, filePprofHeap, fileTrace *os.File
	options                                Options
}

func MakeContext(
	ctx errors.Context,
	options Options,
) (c *Context, err error) {
	c = &Context{
		options: options,
	}

	if options.PProfCPU {
		if c.filePprofCpu, err = files.Create("cpu.pprof"); err != nil {
			err = errors.Wrap(err)
			return
		}

		pprof.StartCPUProfile(c.filePprofCpu)
	}

	if options.Trace {
		if c.fileTrace, err = files.Create("trace"); err != nil {
			err = errors.Wrap(err)
			return
		}

		trace.Start(c.fileTrace)
	}

	if options.GCDisabled {
		debug.SetGCPercent(-1)
	}

	ctx.After(c.Close)

	return
}

func (c *Context) Close() error {
	waitGroupStopOrWrite := errors.MakeWaitGroupParallel()
	multiError := errors.MakeMulti()

	if c.fileTrace != nil {
		waitGroupStopOrWrite.Do(errors.MakeNilFunc(trace.Stop))
	}

	if c.filePprofCpu != nil {
		waitGroupStopOrWrite.Do(errors.MakeNilFunc(pprof.StopCPUProfile))
	}

	if c.options.PProfHeap {
		{
			var err error

			if c.filePprofHeap, err = files.Create("heap.pprof"); err != nil {
				multiError.Add(errors.Wrap(err))
			}
		}

		waitGroupStopOrWrite.Do(func() error {
			return pprof.WriteHeapProfile(c.filePprofHeap)
		})
	}

	if err := waitGroupStopOrWrite.GetError(); err != nil {
		multiError.Add(errors.Wrap(err))
	}

	waitGroupClose := errors.MakeWaitGroupParallel()

	if c.fileTrace != nil {
		waitGroupClose.Do(c.fileTrace.Close)
	}

	if c.filePprofCpu != nil {
		waitGroupClose.Do(c.filePprofCpu.Close)
	}

	if c.options.PProfHeap {
		waitGroupClose.Do(c.filePprofHeap.Close)
	}

	if err := waitGroupClose.GetError(); err != nil {
		multiError.Add(errors.Wrap(err))
	}

	if multiError.Len() > 0 {
		return multiError
	}

	return nil
}
