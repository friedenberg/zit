package zettel_formats

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/charlie/line_format"
	"github.com/friedenberg/zit/charlie/zk_types"
	"github.com/friedenberg/zit/delta/etikett"
	"github.com/friedenberg/zit/echo/zettel"
)

type Objekte struct{}

func (f Objekte) WriteTo(z zettel.Zettel, out1 io.Writer) (n int64, err error) {
	w := line_format.NewWriter()

	w.WriteFormat("%s %s", zk_types.TypeAkte, z.Akte)
	w.WriteFormat("%s %s", zk_types.TypeAkteTyp, z.AkteExt)
	w.WriteFormat("%s %s", zk_types.TypeBezeichnung, z.Bezeichnung)

	for _, e := range z.Etiketten {
		w.WriteFormat("%s %s", zk_types.TypeEtikett, e)
	}

	n, err = w.WriteTo(out1)

	if err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (f *Objekte) ReadFrom(z *zettel.Zettel, in io.Reader) (n int64, err error) {
	z.Etiketten = etikett.MakeSet()

	r := bufio.NewReader(in)

	l := 0

	for {
		var line string
		line, err = r.ReadString('\n')

		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			return
		}

		n += int64(len(line))

		loc := strings.Index(line, " ")

		if loc == -1 {
			err = errors.Errorf("expected at least one space, but found none: %s", line)
			return
		}

		var t zk_types.Type

		if err = t.Set(line[:loc]); err != nil {
			err = errors.Errorf("%s: %s", err, line[:loc])
			return
		}

		v := line[loc+1:]

		switch l {
		case 0:
			if t != zk_types.TypeAkte {
				err = errors.Errorf("expected type %s, but got %s", zk_types.TypeAkte, t)
				return
			}

			if err = z.Akte.Set(v); err != nil {
				err = errors.Error(err)
				return
			}

		case 1:
			if t != zk_types.TypeAkteTyp {
				err = errors.Errorf("expected type %s, but got %s: %s", zk_types.TypeAkteTyp, t, line)
				return
			}

			//TODO: switch to AkteTyp
			if err = z.AkteExt.Set(v); err != nil {
				err = errors.Error(err)
				return
			}

		case 2:
			if t != zk_types.TypeBezeichnung {
				err = errors.Errorf("expected type %s, but got %s", zk_types.TypeBezeichnung, t)
				return
			}

			if err = z.Bezeichnung.Set(v); err != nil {
				err = errors.Error(err)
				return
			}

		default:
			if t != zk_types.TypeEtikett {
				err = errors.Errorf("expected type %s, but got %s", zk_types.TypeEtikett, t)
				return
			}

			if err = z.Etiketten.AddString(v); err != nil {
				err = errors.Error(err)
				return
			}

		}

		l += 1
	}

	if l < 3 {
		err = errors.Errorf("expected at least 3 objekte refs, but got %d", l)
		return
	}

	return
}
