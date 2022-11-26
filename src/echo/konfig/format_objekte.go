package konfig

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/line_format"
	"github.com/friedenberg/zit/src/delta/metadatei_io"
)

type FormatObjekte struct {
	arf metadatei_io.AkteIOFactory
}

func MakeFormatObjekte(arf metadatei_io.AkteIOFactory) *FormatObjekte {
	return &FormatObjekte{
		arf: arf,
	}
}

func (f FormatObjekte) ReadFormat(r1 io.Reader, t *Stored) (n int64, err error) {
	r := bufio.NewReader(r1)

	for {
		var lineOriginal string
		lineOriginal, err = r.ReadString('\n')

		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return
		}

		// line := strings.TrimSpace(lineOriginal)
		line := lineOriginal

		n += int64(len(lineOriginal))

		loc := strings.Index(line, " ")

		if line == "" {
			//TODO this should be cleaned up
		}

		var g gattung.Gattung

		switch {
		case line == "":
			err = errors.Errorf("found empty line: %q", lineOriginal)
			return

		case line != "" && loc == -1:
			if err = g.Set(line[:loc]); err != nil {
				err = errors.Errorf("%s: %s", err, line[:loc])
				return
			}

			if g != gattung.Konfig {
				err = errors.Errorf(
					"expected objekte to have gattung '%s' but got '%s'",
					gattung.Konfig,
					g,
				)

				return
			}

			continue

		case lineOriginal == "\n" && loc == -1:
			continue

		case loc == -1:
			err = errors.Errorf("expected at least one space, but found none: %q", lineOriginal)
			return
		}

		if err = g.Set(line[:loc]); err != nil {
			err = errors.Errorf("%s: %s", err, line[:loc])
			return
		}

		v := line[loc+1:]

		switch g {
		case gattung.Akte:
			if err = t.Sha.Set(v); err != nil {
				err = errors.Wrap(err)
				return
			}

		case gattung.Konfig:

		default:
			err = errors.Errorf("unsupported gattung: %s", g)
			return
		}
	}

	return
}

func (f FormatObjekte) WriteFormat(w1 io.Writer, t *Stored) (n int64, err error) {
	w := line_format.NewWriter()

	w.WriteFormat("%s", gattung.Typ)
	w.WriteFormat("%s %s", gattung.Akte, t.Objekte.Sha)

	if n, err = w.WriteTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
