package sku_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type StringFormatWriterCheckedOutLike = schnittstellen.StringFormatWriter[sku.CheckedOutLike]

type cliCheckedOutLike struct {
	externalWriters map[string]StringFormatWriterCheckedOutLike
}

func MakeCliCheckedOutLikeFormat(
	externalWriters map[string]StringFormatWriterCheckedOutLike,
) *cliCheckedOutLike {
	return &cliCheckedOutLike{
		externalWriters: externalWriters,
	}
}

func (f *cliCheckedOutLike) WriteStringFormat(
	sw schnittstellen.WriterAndStringWriter,
	co sku.CheckedOutLike,
) (n int64, err error) {
	kid := co.GetKasten().GetKastenString()
	sfw, ok := f.externalWriters[kid]

	if !ok {
		err = errors.Errorf("unsupported check out type: %T", co)
		return
	}

	if n, err = sfw.WriteStringFormat(sw, co); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
