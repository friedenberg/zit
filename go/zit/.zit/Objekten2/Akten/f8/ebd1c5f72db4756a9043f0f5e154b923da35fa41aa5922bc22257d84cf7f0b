package sha_probe_index

import (
	"bytes"
	"slices"

	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type addedMap map[sha.Bytes]*row

func (a addedMap) ToSlice() []*row {
	out := make([]*row, 0, len(a))

	for _, r := range a {
		out = append(out, r)
	}

	slices.SortFunc(out, func(x, y *row) int {
		return bytes.Compare(x.left.GetShaBytes(), y.left.GetShaBytes())
	})

	return out
}
