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
		if c.filePprofCpu, err = files.Create("build/cpu.pprof"); err != nil {
			err = errors.Wrap(err)
			return
		}

		pprof.StartCPUProfile(c.filePprofCpu)
	}

	if options.Trace {
		if c.fileTrace, err = files.Create("build/trace"); err != nil {
			err = errors.Wrap(err)
			return
		}

		trace.Start(c.fileTrace)
	}

	if options.GCDisabled {
		debug.SetGCPercent(-1)
	}

	go func() {
		<-ctx.Done()
		c.Close()
	}()

	return
}

func (c *Context) Close() (err error) {
	if c.fileTrace != nil {
		trace.Stop()

		if err = c.fileTrace.Close(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if c.filePprofCpu != nil {
		pprof.StopCPUProfile()

		if err = c.filePprofCpu.Close(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if c.options.PProfHeap {
		if c.filePprofHeap, err = files.Create("build/heap.pprof"); err != nil {
			err = errors.Wrap(err)
			return
		}

		pprof.WriteHeapProfile(c.filePprofHeap)

		if err = c.filePprofHeap.Close(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
