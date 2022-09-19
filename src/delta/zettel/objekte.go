package zettel

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/line_format"
	"github.com/friedenberg/zit/src/bravo/zk_types"
	"github.com/friedenberg/zit/src/charlie/etikett"
)

type Objekte struct {
	IgnoreTypErrors bool
}

func (f Objekte) WriteTo(z Zettel, out1 io.Writer) (n int64, err error) {
	w := line_format.NewWriter()

	w.WriteFormat("%s %s", zk_types.TypeAkte, z.Akte)
	w.WriteFormat("%s %s", zk_types.TypeAkteTyp, z.Typ)
	w.WriteFormat("%s %s", zk_types.TypeBezeichnung, z.Bezeichnung)

	for _, e := range z.Etiketten.Sorted() {
		w.WriteFormat("%s %s", zk_types.TypeEtikett, e)
	}

	n, err = w.WriteTo(out1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *Objekte) ReadFrom(z *Zettel, in io.Reader) (n int64, err error) {
	z.Etiketten = etikett.MakeSet()

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

		var t zk_types.Type

		if err = t.Set(line[:loc]); err != nil {
			err = errors.Errorf("%s: %s", err, line[:loc])
			return
		}

		v := line[loc+1:]

		switch t {
		case zk_types.TypeAkte:
			if err = z.Akte.Set(v); err != nil {
				err = errors.Wrap(err)
				return
			}

		case zk_types.TypeAkteTyp:
			if f.IgnoreTypErrors {
				z.Typ.Etikett.Value = strings.TrimSpace(v)
			} else {
				if err = z.Typ.Set(v); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

		case zk_types.TypeBezeichnung:
			if err = z.Bezeichnung.Set(v); err != nil {
				err = errors.Wrap(err)
				return
			}

		case zk_types.TypeEtikett:
			if err = z.Etiketten.AddString(v); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}
