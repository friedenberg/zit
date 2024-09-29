package object_metadata

import (
	"bufio"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
)

type Reader struct {
	state           readerState
	RequireMetadata bool // TODO-P4 add delimiter
	Metadata, Blob  io.ReaderFrom
}

// TODO-P4 add constructors and remove public fields
func (mr *Reader) ReadFrom(r io.Reader) (n int64, err error) {
	var n1 int64
	n1, err = mr.ReadMetadataFrom(&r)
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

func (mr *Reader) ReadMetadataFrom(r *io.Reader) (n int64, err error) {
	br := bufio.NewReader(*r)

	if mr.RequireMetadata && mr.Metadata == nil {
		err = errors.Errorf("metadata reader is nil")
		return
	}

	if mr.Blob == nil {
		err = errors.Errorf("blob reader is nil")
		return
	}

	var object_metadata ohio.PipedReaderFrom

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
			case mr.RequireMetadata && line != Boundary:
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

			object_metadata = ohio.MakePipedReaderFrom(mr.Metadata)

		case readerStateFirstBoundary:
			if line == Boundary {
				if _, err = object_metadata.Close(); err != nil {
					err = errors.Wrapf(err, "metadata read failed")
					return
				}

				mr.state += 1
				break
			}

			if _, err = object_metadata.Write([]byte(rawLine)); err != nil {
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
	blob := ohio.MakePipedReaderFrom(mr.Blob)

	var n1 int64
	n1, err = br.WriteTo(blob)
	n += n1

	if err != nil {
		err = errors.Wrapf(err, "blob write failed")
		return
	}

	if _, err = blob.Close(); err != nil {
		err = errors.Wrapf(err, "blob read failed")
		return
	}

	return
}
