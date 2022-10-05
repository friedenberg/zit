package debug

import (
	"os"
	"runtime/pprof"
	"runtime/trace"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
)

type Context struct {
	filePprof, fileTrace *os.File
}

func MakeContext(options Options) (c *Context, err error) {
	c = &Context{}

	if c.filePprof, err = files.Create("build/cpu1.pprof"); err != nil {
		err = errors.Wrap(err)
		return
	}

	pprof.StartCPUProfile(c.filePprof)

	if c.fileTrace, err = files.Create("build/trace"); err != nil {
		err = errors.Wrap(err)
		return
	}

	trace.Start(c.fileTrace)

	// debug.SetGCPercent(-1)

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

	if c.filePprof != nil {
		pprof.StopCPUProfile()

		if err = c.filePprof.Close(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
