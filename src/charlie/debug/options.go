package debug

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/string_builder_joined"
)

type Options struct {
	Trace, PProfCPU, PProfHeap, GCDisabled bool
}

func (o Options) String() string {
	sb := string_builder_joined.Make(",")

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
