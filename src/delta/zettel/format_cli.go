package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/bezeichnung"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/typ"
)

func MakeCliFormat(
	af collections.WriterFuncFormat[sha.Sha],
	bf collections.WriterFuncFormat[bezeichnung.Bezeichnung],
	ef collections.WriterFuncFormat[etikett.Set],
	tf collections.WriterFuncFormat[typ.Typ],
) collections.WriterFuncFormat[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
		//TODO
		return
	}
}
