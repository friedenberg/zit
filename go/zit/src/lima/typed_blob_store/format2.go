package typed_blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type format2[
	O interfaces.Blob[O],
] struct {
	interfaces.DecoderFrom[O]
	interfaces.SavedBlobFormatter
	interfaces.EncoderTo[O]
}

func MakeBlobFormat2[
	O interfaces.Blob[O],
](
	decoder interfaces.DecoderFrom[O],
	encoder interfaces.EncoderTo[O],
	blobReader interfaces.BlobReader,
) interfaces.Format[O] {
	return format2[O]{
		DecoderFrom:        decoder,
		EncoderTo:          encoder,
		SavedBlobFormatter: MakeSavedBlobFormatter(blobReader),
	}
}
