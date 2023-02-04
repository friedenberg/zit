package zettel

import (
	"encoding/gob"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/kennung_index"
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
	if err = w.enc.Encode(z); err != nil {
		err = errors.Wrap(err)
    errors.Err().Printf("decode error: %s", err)
		errors.TodoP0("make sure Verzeichnisse flush prints errors")
		return
	}

	return
}
