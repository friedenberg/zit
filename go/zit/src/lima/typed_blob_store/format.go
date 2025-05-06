package typed_blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type format[
	O any,
	OPtr interfaces.Ptr[O],
] struct {
	interfaces.DecoderFrom[OPtr]
	interfaces.SavedBlobFormatter
	interfaces.EncoderTo[OPtr]
}

func MakeBlobFormat[
	O any,
	OPtr interfaces.Ptr[O],
](
	decoder interfaces.DecoderFrom[OPtr],
	encoder interfaces.EncoderTo[OPtr],
	blobReader interfaces.BlobReader,
) Format[O, OPtr] {
	return format[O, OPtr]{
		DecoderFrom:        decoder,
		EncoderTo:          encoder,
		SavedBlobFormatter: MakeSavedBlobFormatter(blobReader),
	}
}

func (af format[O, OPtr]) EncodeTo(
	object OPtr,
	writer io.Writer,
) (n int64, err error) {
	if af.EncoderTo == nil {
		err = errors.ErrorWithStackf("no ParsedBlobFormatter")
	} else {
		n, err = af.EncoderTo.EncodeTo(object, writer)
	}

	return
}
