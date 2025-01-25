package triple_hyphen_io

import (
	"bufio"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
)

type Coder[O any] struct {
	RequireMetadata bool
	Metadata, Blob  interfaces.Coder[O]
}

func (mr *Coder[O]) DecodeFrom(object O, r io.Reader) (n int64, err error) {
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

func (mr *Coder[O]) readMetadataFrom(
	object O,
	r *io.Reader,
) (n int64, err error) {
	var state readerState
	br := bufio.NewReader(*r)

	if mr.RequireMetadata && mr.Metadata == nil {
		err = errors.Errorf("metadata reader is nil")
		return
	}

	if mr.Blob == nil {
		err = errors.Errorf("blob reader is nil")
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
				err = errors.Errorf("expected %q but got %q", Boundary, line)
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
			err = errors.Errorf("impossible state %d", state)
			return
		}
	}

	return
}

func (w1 Coder[O]) EncodeTo(
	object O,
	w2 io.Writer,
) (n int64, err error) {
	w := bufio.NewWriter(w2)
	defer errors.DeferredFlusher(&err, w)

	var n1 int64
	var n2 int

	hasMetadataContent := true

	if mwt, ok := w1.Metadata.(MetadataWriterTo); ok {
		hasMetadataContent = mwt.HasMetadataContent()
	}

	if w1.Metadata != nil && hasMetadataContent {
		n2, err = w.WriteString(Boundary + "\n")
		n += int64(n2)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n1, err = w1.Metadata.EncodeTo(object, w)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		w.WriteString(Boundary + "\n")
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		if w1.Blob != nil {
			w.WriteString("\n")
			n += n1

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	if w1.Blob != nil {
		n1, err = w1.Blob.EncodeTo(object, w)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
