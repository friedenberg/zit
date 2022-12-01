package zettel

import (
	"bufio"
	"crypto/sha256"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/line_format"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
)

func (z Zettel) ObjekteSha() (s sha.Sha, err error) {
	hash := sha256.New()

	o := FormatObjekte{}

	c := FormatContextWrite{
		Zettel: z,
		Out:    hash,
	}

	if _, err = o.WriteTo(c); err != nil {
		err = errors.Wrap(err)
		return
	}

	s = sha.FromHash(hash)

	return
}

type FormatObjekte struct {
	IgnoreTypErrors bool
}

func (f FormatObjekte) WriteTo(c FormatContextWrite) (n int64, err error) {
	z := c.Zettel
	w := line_format.NewWriter()

	w.WriteFormat("%s %s", gattung.Akte, z.Akte)
	w.WriteFormat("%s %s", gattung.Typ, z.Typ)
	w.WriteFormat("%s %s", gattung.Bezeichnung, z.Bezeichnung)

	for _, e := range z.Etiketten.Sorted() {
		w.WriteFormat("%s %s", gattung.Etikett, e)
	}

	n, err = w.WriteTo(c.Out)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *FormatObjekte) ReadFrom(c *FormatContextRead) (n int64, err error) {
	etiketten := kennung.MakeEtikettMutableSet()

	var z *Zettel
	z = &c.Zettel

	r := bufio.NewReader(c.In)

	bezLineReader := z.Bezeichnung.Set

	if f.IgnoreTypErrors {
		bezLineReader = format.MakeLineReaderIgnoreErrors(bezLineReader)
	}

	if n, err = format.ReadLines(
		r,
		format.MakeLineReaderKeyValue(gattung.Akte.String(), z.Akte.Set),
		format.MakeLineReaderKeyValue(gattung.Typ.String(), z.Typ.Set),
		format.MakeLineReaderKeyValue(gattung.Bezeichnung.String(), bezLineReader),
		format.MakeLineReaderKeyValue(gattung.Etikett.String(), etiketten.AddString),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	z.Etiketten = etiketten.Copy()

	return
}
