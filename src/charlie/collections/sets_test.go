package collections

import (
	"github.com/friedenberg/zit/src/bravo/values"
)

func makeStringValues(vs ...string) (out []values.String) {
	out = make([]values.String, len(vs))

	for i, v := range vs {
		out[i] = values.MakeString(v)
	}

	return
}
