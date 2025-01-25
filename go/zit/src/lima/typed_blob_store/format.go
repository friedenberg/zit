package typed_blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type format[
	O interfaces.Blob[O],
	OPtr interfaces.BlobPtr[O],
] struct {
	interfaces.DecoderFrom[OPtr]
	interfaces.SavedBlobFormatter
	interfaces.EncoderTo[OPtr]
}

func MakeBlobFormat[
	O interfaces.Blob[O],
	OPtr interfaces.BlobPtr[O],
](
	decoder interfaces.DecoderFrom[OPtr],
	encoder interfaces.EncoderTo[OPtr],
	arf interfaces.BlobReader,
) Format[O, OPtr] {
	return format[O, OPtr]{
		DecoderFrom:        decoder,
		EncoderTo:          encoder,
		SavedBlobFormatter: MakeSavedBlobFormatter(arf),
	}
}

func (af format[O, OPtr]) EncodeTo(
	object OPtr,
	writer io.Writer,
) (n int64, err error) {
	if af.EncoderTo == nil {
		err = errors.Errorf("no ParsedBlobFormatter")
	} else {
		n, err = af.EncoderTo.EncodeTo(object, writer)
	}

	return
}
