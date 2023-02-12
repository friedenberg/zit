package metadatei_io

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Reader struct {
	RequireMetadatei bool
	Metadatei, Akte  io.ReaderFrom
}

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

	var metadatei, akte pipedReaderFrom
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

			metadatei = makePipedReaderFrom(mr.Metadatei)

		case readerStateFirstBoundary:
			if line == Boundary {
				msg := metadatei.Close()

				if msg.err != nil {
					err = errors.Wrapf(msg.err, "metadatei read failed")
					return
				}

				state += 1
				break
			}

			if _, err = metadatei.PipeWriter.Write([]byte(rawLine)); err != nil {
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

	akte = makePipedReaderFrom(mr.Akte)

	var n1 int64
	n1, err = r.WriteTo(akte.PipeWriter)
	n += n1

	msg := akte.Close()

	if msg.err != nil {
		err = errors.Wrapf(msg.err, "akte read failed")
		return
	}

	// TODO handle errors

	return
}
