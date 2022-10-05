package debug

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/string_joined_builder"
)

type Options struct {
	Trace, PProf, GCDisabled bool
}

func (o Options) String() string {
	sb := string_joined_builder.Make(",")

	if o.GCDisabled {
		sb.WriteString("gc_disabled")
	}

	if o.PProf {
		sb.WriteString("pprof")
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

		case "pprof":
			o.PProf = true

		case "trace":
			o.Trace = true

		case "true":
      fallthrough

		case "all":
			o.GCDisabled = true
			o.PProf = true
			o.Trace = true

		default:
			err = errors.Errorf("unsupported debug option: %s", p)
			return
		}
	}

	return
}
