package metadatei

import (
	"bufio"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
)

type Reader struct {
	state            readerState
	RequireMetadatei bool // TODO-P4 add delimiter
	Metadatei, Akte  io.ReaderFrom
}

// TODO-P4 add constructors and remove public fields
func (mr *Reader) ReadFrom(r io.Reader) (n int64, err error) {
	var n1 int64
	n1, err = mr.ReadMetadateiFrom(&r)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = mr.ReadBlobFrom(r)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (mr *Reader) ReadMetadateiFrom(r *io.Reader) (n int64, err error) {
	br := bufio.NewReader(*r)

	if mr.RequireMetadatei && mr.Metadatei == nil {
		err = errors.Errorf("metadatei reader is nil")
		return
	}

	if mr.Akte == nil {
		err = errors.Errorf("akte reader is nil")
		return
	}

	var metadatei ohio.PipedReaderFrom

	isEOF := false

LINE_READ_LOOP:
	for !isEOF {
		var rawLine, line string

		rawLine, err = br.ReadString('\n')
		n += int64(len(rawLine))

		if err != nil && err != io.EOF {
			err = errors.Wrap(err)
			return
		}

		if errors.IsEOF(err) {
			err = nil
			isEOF = true
		}

		line = strings.TrimSuffix(rawLine, "\n")

		switch mr.state {
		case readerStateEmpty:
			switch {
			case mr.RequireMetadatei && line != Boundary:
				err = errors.Errorf("expected %q but got %q", Boundary, line)
				return

			case line != Boundary:
				*r = io.MultiReader(
					strings.NewReader(rawLine),
					br,
				)

				break LINE_READ_LOOP
			}

			mr.state += 1

			metadatei = ohio.MakePipedReaderFrom(mr.Metadatei)

		case readerStateFirstBoundary:
			if line == Boundary {
				if _, err = metadatei.Close(); err != nil {
					err = errors.Wrapf(err, "metadatei read failed")
					return
				}

				mr.state += 1
				break
			}

			if _, err = metadatei.Write([]byte(rawLine)); err != nil {
				err = errors.Wrap(err)
				return
			}

		case readerStateSecondBoundary:
			*r = br
			break LINE_READ_LOOP

		default:
			err = errors.Errorf("impossible state %d", mr.state)
			return
		}
	}

	return
}

func (mr *Reader) ReadBlobFrom(r io.Reader) (n int64, err error) {
	br := bufio.NewReader(r)
	akte := ohio.MakePipedReaderFrom(mr.Akte)

	var n1 int64
	n1, err = br.WriteTo(akte)
	n += n1

	if err != nil {
		err = errors.Wrapf(err, "akte write failed")
		return
	}

	if _, err = akte.Close(); err != nil {
		err = errors.Wrapf(err, "akte read failed")
		return
	}

	return
}
