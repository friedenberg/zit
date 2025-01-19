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
	interfaces.ParsedBlobFormatter[O]
}

func MakeBlobFormat2[
	O interfaces.Blob[O],
](
	parser interfaces.Parser[O],
	formatter interfaces.ParsedBlobFormatter[O],
	arf interfaces.BlobReader,
) interfaces.Format[O] {
	return format2[O]{
		Parser:              parser,
		ParsedBlobFormatter: formatter,
		SavedBlobFormatter:  MakeSavedBlobFormatter(arf),
	}
}

func (af format2[O]) FormatParsedBlob(
	w io.Writer,
	e O,
) (n int64, err error) {
	if af.ParsedBlobFormatter == nil {
		err = errors.Errorf("no ParsedBlobFormatter")
	} else {
		n, err = af.ParsedBlobFormatter.FormatParsedBlob(w, e)
	}

	return
}
