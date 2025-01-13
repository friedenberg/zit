package store_fs

import (
	"maps"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

type fsItemData struct {
	interfaces.MutableSetLike[*sku.FSItem]
	shas map[sha.Bytes]interfaces.MutableSetLike[*sku.FSItem]
}

func makeFSItemData() fsItemData {
	return fsItemData{
		MutableSetLike: collections_value.MakeMutableValueSet[*sku.FSItem](nil),
		shas:           make(map[sha.Bytes]interfaces.MutableSetLike[*sku.FSItem]),
	}
}

func (src *fsItemData) Clone() (dst fsItemData) {
	dst.MutableSetLike = src.MutableSetLike.CloneMutableSetLike()
	dst.shas = maps.Clone(src.shas)
	return
}

func (data *fsItemData) ConsolidateDuplicateBlobs() (err error) {
	replacement := collections_value.MakeMutableValueSet[*sku.FSItem](nil)

	for _, fds := range data.shas {
		if fds.Len() == 1 {
			replacement.Add(fds.Any())
		}

		sorted := quiter.ElementsSorted(
			fds,
			func(a, b *sku.FSItem) bool {
				return a.ExternalObjectId.String() < b.ExternalObjectId.String()
			},
		)

		top := sorted[0]

		for _, other := range sorted[1:] {
			other.MutableSetLike.Each(top.MutableSetLike.Add)
		}

		replacement.Add(top)
	}

	// TODO make less leaky
	data.MutableSetLike = replacement

	return
}
