package zettel

import (
	"encoding/gob"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/foxtrot/kennung_index"
	"github.com/friedenberg/zit/src/india/konfig"
)

func init() {
	errors.TodoP0("add hidden")
}

type Verzeichnisse struct {
	wasPopulated bool
	Etiketten    kennung_index.EtikettenVerzeichnisse
	Typ          kennung_index.TypVerzeichnisse
	// Hidden bool
}

func (z *Verzeichnisse) ResetWithObjekte(z1 Objekte) {
	z.wasPopulated = true

	z.Etiketten.ResetWithEtikettSet(z1.Etiketten)
	z.Typ.ResetWithTyp(z1.Typ)
}

func (z *Verzeichnisse) Reset() {
	z.wasPopulated = false

	z.Etiketten.Reset()
	z.Typ.Reset()
}

func (z *Verzeichnisse) ResetWith(z1 Verzeichnisse) {
	z.wasPopulated = true

	z.Etiketten.ResetWith(z1.Etiketten)
	z.Typ.ResetWith(z1.Typ)
}

// TODO-P1 remove and make transform func in collections
type writerGobEncoder struct {
	enc *gob.Encoder
}

func MakeWriterGobEncoder(w io.Writer) writerGobEncoder {
	return writerGobEncoder{
		enc: gob.NewEncoder(w),
	}
}

func (w writerGobEncoder) WriteZettelVerzeichnisse(z *Transacted) (err error) {
	return w.enc.Encode(z)
}

func MakeWriterKonfig(
	k konfig.Compiled,
) collections.WriterFunc[*Transacted] {
	errors.TodoP3("add efficient parsing of hiding tags")

	if k.IncludeHidden {
		return collections.MakeWriterNoop[*Transacted]()
	}

	return func(z *Transacted) (err error) {
		for _, p := range z.Verzeichnisse.Etiketten.Sorted {
			for _, t := range k.EtikettenHidden {
				if strings.HasPrefix(p, t) {
					err = collections.MakeErrStopIteration()
					return
				}
			}
		}

		return
	}
}
