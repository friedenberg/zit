package objekte_format

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/line_format"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/metadatei_io"
)

type Objekte interface {
	Gattung() gattung.Gattung
}

type Objekte2 interface {
	Objekte

	AkteSha() sha.Sha
}

type ObjektePtr[T any] interface {
	*T
	collections.Equatable[T]
	collections.Resetable[T]
}

type Objekte2Ptr[T any] interface {
	*T
	collections.Equatable[T]
	collections.Resetable[T]

	SetAkteSha(sha.Sha)
}

type Stored2 interface {
	Gattung() gattung.Gattung

	AkteSha() sha.Sha
	SetAkteSha(sha.Sha)

	SetObjekteSha(metadatei_io.AkteReaderFactory, string) error
	ObjekteSha() sha.Sha
}

type Format struct {
	arf metadatei_io.AkteIOFactory
}

func MakeFormat(arf metadatei_io.AkteIOFactory) *Format {
	return &Format{
		arf: arf,
	}
}

func (f Format) ReadFormat(
	r1 io.Reader,
	o Stored2,
) (n int64, err error) {
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

			if g != o.Gattung() {
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
			var sh sha.Sha

			if err = sh.Set(v); err != nil {
				err = errors.Wrap(err)
				return
			}

			o.SetAkteSha(sh)

		default:
			err = errors.Errorf("unsupported gattung: %s", g)
			return
		}
	}

	return
}

func (f Format) WriteFormat(
	w1 io.Writer,
	o Stored2,
) (n int64, err error) {
	w := line_format.NewWriter()

	w.WriteFormat("%s", o.Gattung())
	w.WriteFormat("%s %s", gattung.Akte, o.AkteSha())

	if n, err = w.WriteTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
