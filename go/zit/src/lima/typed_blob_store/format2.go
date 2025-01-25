package typed_blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type format2[
	O interfaces.Blob[O],
] struct {
	interfaces.Parser[O]
	interfaces.ParseSaver[O]
	interfaces.SavedBlobFormatter
	interfaces.EncoderTo[O]
}

func MakeBlobFormat2[
	O interfaces.Blob[O],
](
	parser interfaces.Parser[O],
	formatter interfaces.EncoderTo[O],
	arf interfaces.BlobReader,
) interfaces.Format[O] {
	return format2[O]{
		Parser:             parser,
		EncoderTo:          formatter,
		SavedBlobFormatter: MakeSavedBlobFormatter(arf),
	}
}

func (af format2[O]) FormatParsedBlob(
	w io.Writer,
	e O,
) (n int64, err error) {
	if af.EncoderTo == nil {
		err = errors.Errorf("no ParsedBlobFormatter")
	} else {
		n, err = af.EncoderTo.EncodeTo(e, w)
	}

	return
}
