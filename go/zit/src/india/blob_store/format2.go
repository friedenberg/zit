package blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type format2[
	O interfaces.Blob[O],
] struct {
	Parser2[O]
	ParseSaver2[O]
	SavedBlobFormatter
	ParsedBlobFormatter2[O]
}

func MakeBlobFormat2[
	O interfaces.Blob[O],
](
	parser Parser2[O],
	formatter ParsedBlobFormatter2[O],
	arf interfaces.BlobReaderFactory,
) Format2[O] {
	return format2[O]{
		Parser2:              parser,
		ParsedBlobFormatter2: formatter,
		SavedBlobFormatter:   MakeSavedBlobFormatter(arf),
	}
}

func (af format2[O]) FormatParsedBlob(
	w io.Writer,
	e O,
) (n int64, err error) {
	if af.ParsedBlobFormatter2 == nil {
		err = errors.Errorf("no ParsedBlobFormatter")
	} else {
		n, err = af.ParsedBlobFormatter2.FormatParsedBlob(w, e)
	}

	return
}
