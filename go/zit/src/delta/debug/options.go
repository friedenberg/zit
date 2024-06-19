package debug

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/string_builder_joined"
)

type Options struct {
	Trace            bool
	PProfCPU         bool
	PProfHeap        bool
	GCDisabled       bool
	NoTempDirCleanup bool
	DryRun           bool
}

func (o Options) String() string {
	sb := string_builder_joined.Make(",")

	if o.NoTempDirCleanup {
		sb.WriteString("no-tempdir-cleanup")
	}

	if o.GCDisabled {
		sb.WriteString("gc_disabled")
	}

	if o.PProfCPU {
		sb.WriteString("pprof_cpu")
	}

	if o.PProfHeap {
		sb.WriteString("pprof_heap")
	}

	if o.Trace {
		sb.WriteString("trace")
	}

	if o.DryRun {
		sb.WriteString("dry-run")
	}

	return sb.String()
}

func (o *Options) Set(v string) (err error) {
	parts := strings.Split(v, ",")

	if len(parts) == 0 {
		parts = []string{"all"}
	}

	for _, p := range parts {
		switch strings.ToLower(p) {
		case "false":

		case "gc_disabled":
			o.GCDisabled = true

		case "pprof_cpu":
			o.PProfCPU = true

		case "pprof_heap":
			o.PProfHeap = true

		case "trace":
			o.Trace = true

		case "no-tempdir-cleanup":
			o.NoTempDirCleanup = true

		case "dry-run":
			o.DryRun = true

		case "true":
			fallthrough

		case "all":
			o.GCDisabled = true
			o.PProfCPU = true
			o.PProfHeap = true
			o.Trace = true

		default:
			err = errors.Errorf("unsupported debug option: %s", p)
			return
		}
	}

	return
}
