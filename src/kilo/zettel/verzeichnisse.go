package zettel

import (
	"encoding/gob"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/juliett/konfig"
)

func init() {
	errors.TodoP0("add typen expanded sorted")
	errors.TodoP0("add hidden")
}

type Verzeichnisse struct {
	wasPopulated bool
	// Etiketten               tridex.Tridex
	EtikettenExpandedSorted []string
	EtikettenSorted         []string
	// Hidden bool
}

func (z *Verzeichnisse) ResetWithObjekte(z1 *Objekte) {
	if z1 != nil {
		z.wasPopulated = true
		z.EtikettenExpandedSorted = kennung.Expanded(z1.Etiketten).SortedString()
		z.EtikettenSorted = z1.Etiketten.SortedString()
	} else {
		z.wasPopulated = false
		z.EtikettenExpandedSorted = []string{}
		z.EtikettenSorted = []string{}
	}
}

func (z *Verzeichnisse) Reset(z1 *Verzeichnisse) {
	z.ResetWithObjekte(nil)

	if z1 == nil {
		return
	}

	z.EtikettenExpandedSorted = append(
		z.EtikettenExpandedSorted,
		z1.EtikettenExpandedSorted...,
	)

	z.EtikettenSorted = append(
		z.EtikettenSorted,
		z1.EtikettenSorted...,
	)
}

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
		for _, p := range z.Verzeichnisse.EtikettenSorted {
			for _, t := range k.EtikettenHidden {
				if strings.HasPrefix(p, t) {
					err = collections.ErrStopIteration
					return
				}
			}
		}

		return
	}
}
