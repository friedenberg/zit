package zettel

import (
	"bufio"
	"crypto/sha256"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/line_format"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/etikett"
)

func (z Zettel) ObjekteSha() (s sha.Sha, err error) {
	hash := sha256.New()

	o := Objekte{}

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

type Objekte struct {
	IgnoreTypErrors bool
}

func (f Objekte) WriteTo(c FormatContextWrite) (n int64, err error) {
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

func (f *Objekte) ReadFrom(c *FormatContextRead) (n int64, err error) {
	etiketten := etikett.MakeMutableSet()
	in := c.In
	z := c.Zettel

	r := bufio.NewReader(in)

	for {
		var lineOriginal string
		lineOriginal, err = r.ReadString('\n')

		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			return
		}

		// line := strings.TrimSpace(lineOriginal)
		line := lineOriginal

		n += int64(len(lineOriginal))

		loc := strings.Index(line, " ")

		if line == "" {
			//TODO this should be cleaned up
			err = errors.Errorf("found empty line: %q", lineOriginal)
			return
		}

		if loc == -1 {
			if lineOriginal == "\n" {
				continue
			}

			err = errors.Errorf("expected at least one space, but found none: %q", lineOriginal)
			return
		}

		var t gattung.Gattung

		if err = t.Set(line[:loc]); err != nil {
			err = errors.Errorf("%s: %s", err, line[:loc])
			return
		}

		v := line[loc+1:]

		switch t {
		case gattung.Akte:
			if err = z.Akte.Set(v); err != nil {
				err = errors.Wrap(err)
				return
			}

		case gattung.Typ:
			if f.IgnoreTypErrors {
				z.Typ.Etikett.Value = strings.TrimSpace(v)
			} else {
				if err = z.Typ.Set(v); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

		case gattung.Bezeichnung:
			if err = z.Bezeichnung.Set(v); err != nil {
				err = errors.Wrap(err)
				return
			}

		case gattung.Etikett:
			if err = etiketten.AddString(v); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	z.Etiketten = etiketten.Copy()

	return
}
