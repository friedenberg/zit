package triple_hyphen_io

import (
	"bufio"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
)

type Decoder[O any] struct {
	RequireMetadata bool
	Metadata, Blob  interfaces.DecoderFrom[O]
}

func (mr *Decoder[O]) DecodeFrom(object O, r io.Reader) (n int64, err error) {
	var n1 int64
	n1, err = mr.readMetadataFrom(object, &r)
	n += n1

	if err != nil {
		err = errors.Wrapf(err, "metadata read failed")
		return
	}

	n1, err = mr.Blob.DecodeFrom(object, bufio.NewReader(r))
	n += n1

	if err != nil {
		err = errors.Wrapf(err, "blob read failed")
		return
	}

	return
}

func (mr *Decoder[O]) readMetadataFrom(
	object O,
	r *io.Reader,
) (n int64, err error) {
	var state readerState
	br := bufio.NewReader(*r)

	if mr.RequireMetadata && mr.Metadata == nil {
		err = errors.ErrorWithStackf("metadata reader is nil")
		return
	}

	if mr.Blob == nil {
		err = errors.ErrorWithStackf("blob reader is nil")
		return
	}

	var metadataPipe ohio.PipedReader

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

		switch state {
		case readerStateEmpty:
			switch {
			case mr.RequireMetadata && line != Boundary:
				err = errors.ErrorWithStackf("expected %q but got %q", Boundary, line)
				return

			case line != Boundary:
				*r = io.MultiReader(
					strings.NewReader(rawLine),
					br,
				)

				break LINE_READ_LOOP
			}

			state += 1

			metadataPipe = ohio.MakePipedDecoder(object, mr.Metadata)

		case readerStateFirstBoundary:
			if line == Boundary {
				if _, err = metadataPipe.Close(); err != nil {
					err = errors.Wrapf(err, "metadata read failed")
					return
				}

				state += 1
				break
			}

			if _, err = metadataPipe.Write([]byte(rawLine)); err != nil {
				err = errors.Wrap(err)
				return
			}

		case readerStateSecondBoundary:
			*r = br
			break LINE_READ_LOOP

		default:
			err = errors.ErrorWithStackf("impossible state %d", state)
			return
		}
	}

	return
}
