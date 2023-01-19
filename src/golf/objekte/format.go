package objekte

import (
	"bufio"
	"crypto/sha256"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/format"
)

type Format[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
] struct{}

func MakeFormat[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
]() *Format[T, T1] {
	return &Format[T, T1]{}
}

func (f Format[T, T1]) Parse(
	r1 io.Reader,
	o T1,
) (n int64, err error) {
	r := bufio.NewReader(r1)

	isEOF := false

	for {
		if isEOF {
			break
		}

		var lineOriginal string
		lineOriginal, err = r.ReadString('\n')

		switch {
		case errors.IsEOF(err):
			err = nil
			isEOF = true
			break

		case err != nil:
			err = errors.Wrap(err)
			return
		}

		line := strings.TrimSpace(lineOriginal)

		n += int64(len(lineOriginal))

		loc := strings.Index(line, " ")

		if line == "" {
			//TODO this should be cleaned up
		}

		var g gattung.Gattung

		switch {
		case line == "":
			return

		case line != "" && loc == -1:
			if err = g.Set(line); err != nil {
				err = errors.Errorf("%s: %s", err, line)
				return
			}

			if g != o.GetGattung() {
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

func (f Format[T, T1]) Format(
	w1 io.Writer,
	o T1,
) (n int64, err error) {
	hash := sha256.New()
	w2 := io.MultiWriter(w1, hash)

	w := format.NewLineWriter()

	//TODO-P4 use new format methods
	w.WriteFormat("%s", o.GetGattung())
	w.WriteFormat("%s %s", gattung.Akte, o.GetAkteSha())

	if n, err = w.WriteTo(w2); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
