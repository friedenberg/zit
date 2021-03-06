package stored_zettel_formats

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/bravo/line_format"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

type Objekte struct{}

func (z Objekte) WriteTo(sz stored_zettel.Stored, out1 io.Writer) (n int64, err error) {
	w := line_format.NewWriter()

	w.WriteFormat("Mutter %s", sz.Mutter)
	w.WriteFormat("Kinder %s", sz.Kinder)
	w.WriteFormat("Akte %s", sz.Zettel.Akte)
	w.WriteFormat("AkteExt %s", sz.Zettel.AkteExt)
	w.WriteFormat("Bezeichnung %s", sz.Zettel.Bezeichnung)

	for _, e := range sz.Zettel.Etiketten.Sorted() {
		w.WriteFormat("Etikett %s", e)
	}

	n, err = w.WriteTo(out1)

	if err != nil {
		err = _Error(err)
		return
	}

	return
}

func (f *Objekte) ReadFrom(sz *stored_zettel.Stored, in io.Reader) (n int64, err error) {
	sz.Zettel.Etiketten = etikett.MakeSet()

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
			err = _Errorf("expected at least one space, but found none: %s", line)
			return
		}

		var t _Type

		if err = t.Set(line[:loc]); err != nil {
			err = _Errorf("%s: %s", err, line[:loc])
			return
		}

		v := line[loc+1:]

		switch l {
		case 0:
			if t != _TypeMutter {
				err = _Errorf("expected type %s, but got %s", _TypeMutter, t)
				return
			}

			if err = sz.Mutter.Set(v); err != nil {
				err = _Errorf("%s: %s", err, line)
				return
			}

		case 1:
			if t != _TypeKinder {
				err = _Errorf("expected type %s, but got %s", _TypeKinder, t)
				return
			}

			if err = sz.Kinder.Set(v); err != nil {
				err = _Error(err)
				return
			}

		case 2:
			if t != _TypeAkte {
				err = _Errorf("expected type %s, but got %s", _TypeAkte, t)
				return
			}

			if err = sz.Zettel.Akte.Set(v); err != nil {
				err = _Error(err)
				return
			}

		case 3:
			if t != _TypeAkteExt {
				err = _Errorf("expected type %s, but got %s: %s", _TypeAkteExt, t, line)
				return
			}

			if err = sz.Zettel.AkteExt.Set(v); err != nil {
				err = _Error(err)
				return
			}

		case 4:
			if t != _TypeBezeichnung {
				err = _Errorf("expected type %s, but got %s", _TypeBezeichnung, t)
				return
			}

			if err = sz.Zettel.Bezeichnung.Set(v); err != nil {
				err = _Error(err)
				return
			}

		default:
			if t != _TypeEtikett {
				err = _Errorf("expected type %s, but got %s", _TypeEtikett, t)
				return
			}

			if err = sz.Zettel.Etiketten.AddString(v); err != nil {
				err = _Error(err)
				return
			}

		}

		l += 1
	}

	if l < 3 {
		err = _Errorf("expected at least 3 objekte refs, but got %d", l)
		return
	}

	return
}
