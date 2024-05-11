package debug

import (
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"runtime/trace"
	"syscall"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/src/charlie/files"
)

type Context struct {
	filePprofCpu, filePprofHeap, fileTrace *os.File
	options                                Options
}

func MakeContext(options Options) (c *Context, err error) {
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

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT)

	go func() {
		<-ch
		ui.Err().Print("SIGINT")
		c.Close()
		runtime.Goexit()
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
