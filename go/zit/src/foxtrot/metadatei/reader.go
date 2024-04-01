package metadatei

import (
	"bufio"
	"io"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/ohio"
)

type Reader struct {
	// TODO-P4 add delimiter
	RequireMetadatei bool
	Metadatei, Akte  io.ReaderFrom
}

// TODO-P4 add constructors and remove public fields
func (mr *Reader) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := bufio.NewReader(r1)

	if mr.RequireMetadatei && mr.Metadatei == nil {
		err = errors.Errorf("metadatei reader is nil")
		return
	}

	if mr.Akte == nil {
		err = errors.Errorf("akte reader is nil")
		return
	}

	var metadatei, akte ohio.PipedReaderFrom
	var state readerState

	isEOF := false

LINE_READ_LOOP:
	for {
		if isEOF {
			break
		}

		var rawLine, line string

		rawLine, err = r.ReadString('\n')
		n += int64(len(rawLine))

		if err != nil && err != io.EOF {
			err = errors.Wrap(err)
			return
		}

		if errors.IsEOF(err) {
			isEOF = true
		}

		line = strings.TrimSuffix(rawLine, "\n")

		switch state {
		case readerStateEmpty:
			switch {
			case mr.RequireMetadatei && line != Boundary:
				err = errors.Errorf("expected %q but got %q", Boundary, line)
				return

			case line != Boundary:
				r2 := io.MultiReader(
					strings.NewReader(rawLine),
					r,
				)

				r = bufio.NewReader(r2)
				break LINE_READ_LOOP
			}

			state += 1

			metadatei = ohio.MakePipedReaderFrom(mr.Metadatei)

		case readerStateFirstBoundary:
			if line == Boundary {
				_, err = metadatei.Close()

				if err != nil {
					err = errors.Wrapf(err, "metadatei read failed")
					return
				}

				state += 1
				break
			}

			if _, err = metadatei.Write([]byte(rawLine)); err != nil {
				err = errors.Wrap(err)
				return
			}

		case readerStateSecondBoundary:
			break LINE_READ_LOOP

		default:
			err = errors.Errorf("impossible state %d", state)
			return
		}
	}

	akte = ohio.MakePipedReaderFrom(mr.Akte)

	var n1 int64
	n1, err = r.WriteTo(akte)
	n += n1

	_, err = akte.Close()

	if err != nil {
		err = errors.Wrapf(err, "akte read failed")
		return
	}

	// TODO-P2 handle errors

	return
}
