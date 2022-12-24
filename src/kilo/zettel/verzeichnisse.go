package zettel

import (
	"encoding/gob"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/juliett/konfig_compiled"
)

// TODO-P1 merge into Transacted
type Verzeichnisse struct {
	Transacted Transacted
	Verzeichnisse2
}

type PoolVerzeichnisse = collections.Pool[Verzeichnisse]

func MakeVerzeichnisse(z1 *Transacted) (z2 *Verzeichnisse) {
	z2 = &Verzeichnisse{}
	z2.ResetWithTransacted(z1)
	return
}

func (z *Verzeichnisse) ResetWithTransacted(z1 *Transacted) {
	if z1 != nil {
		z.Transacted.Reset(z1)
		z.EtikettenExpandedSorted = kennung.Expanded(z1.Objekte.Etiketten).SortedString()
		z.EtikettenSorted = z1.Objekte.Etiketten.SortedString()
	} else {
		z.Transacted.Reset(nil)
		z.EtikettenExpandedSorted = []string{}
		z.EtikettenSorted = []string{}
	}
}

func (z *Verzeichnisse) Reset(z1 *Verzeichnisse) {
	z.EtikettenExpandedSorted = z.EtikettenExpandedSorted[:0]
	z.EtikettenSorted = z.EtikettenSorted[:0]

	if z1 != nil {
		z.Transacted.Reset(&z1.Transacted)

		z.EtikettenExpandedSorted = append(
			z.EtikettenExpandedSorted,
			z1.EtikettenExpandedSorted...,
		)

		z.EtikettenSorted = append(
			z.EtikettenSorted,
			z1.EtikettenSorted...,
		)
	} else {
		z.Transacted.Reset(nil)
	}
}

func MakeWriterZettelTransacted(
	wf collections.WriterFunc[*Transacted],
) collections.WriterFunc[*Verzeichnisse] {
	return func(z *Verzeichnisse) (err error) {
		if err = wf(&z.Transacted); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

type writerGobEncoder struct {
	enc *gob.Encoder
}

func MakeWriterGobEncoder(w io.Writer) writerGobEncoder {
	return writerGobEncoder{
		enc: gob.NewEncoder(w),
	}
}

func (w writerGobEncoder) WriteZettelVerzeichnisse(z *Verzeichnisse) (err error) {
	return w.enc.Encode(z)
}

// TODO-P3 add efficient parsing of hiding tags
func MakeWriterKonfig(
	k konfig_compiled.Compiled,
) collections.WriterFunc[*Verzeichnisse] {
	if k.IncludeHidden {
		return collections.MakeWriterNoop[*Verzeichnisse]()
	}

	return func(z *Verzeichnisse) (err error) {
		for _, p := range z.EtikettenSorted {
			for _, t := range k.EtikettenHidden {
				if strings.HasPrefix(p, t) {
					err = io.EOF
					return
				}
			}
		}

		return
	}
}
