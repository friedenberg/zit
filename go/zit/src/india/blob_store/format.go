package blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type format[
	O interfaces.Blob[O],
	OPtr interfaces.BlobPtr[O],
] struct {
	Parser[O, OPtr]
	ParseSaver[O, OPtr]
	SavedBlobFormatter
	ParsedBlobFormatter[O, OPtr]
}

func MakeBlobFormat[
	O interfaces.Blob[O],
	OPtr interfaces.BlobPtr[O],
](
	parser Parser[O, OPtr],
	formatter ParsedBlobFormatter[O, OPtr],
	arf interfaces.BlobReaderFactory,
) Format[O, OPtr] {
	return format[O, OPtr]{
		Parser:              parser,
		ParsedBlobFormatter: formatter,
		SavedBlobFormatter:  MakeSavedBlobFormatter(arf),
	}
}

func (af format[O, OPtr]) FormatParsedBlob(
	w io.Writer,
	e OPtr,
) (n int64, err error) {
	if af.ParsedBlobFormatter == nil {
		err = errors.Errorf("no ParsedBlobFormatter")
	} else {
		n, err = af.ParsedBlobFormatter.FormatParsedBlob(w, e)
	}

	return
}
