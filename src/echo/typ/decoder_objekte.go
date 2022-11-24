package typ

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
)

type DecoderObjekte struct {
	in io.Reader
}

func MakeDecoderObjekte(in io.Reader) *DecoderObjekte {
	return &DecoderObjekte{
		in: in,
	}
}

func (f *DecoderObjekte) Decode(t *Named) (n int64, err error) {
	r := bufio.NewReader(f.in)

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

		var g gattung.Gattung

		if err = g.Set(line[:loc]); err != nil {
			err = errors.Errorf("%s: %s", err, line[:loc])
			return
		}

		v := line[loc+1:]

		switch g {
		case gattung.Akte:
			if err = t.Stored.Sha.Set(v); err != nil {
				err = errors.Wrap(err)
				return
			}

		case gattung.Typ:

		default:
			err = errors.Errorf("unsupported gattung: %s", g)
			return
		}
	}

	return
}
