package ohio

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type PipedReader interface {
	Close() (n int64, err error)
	io.Writer
}

type readFromDone struct {
	n   int64
	err error
}

type pipedReaderFrom struct {
	*io.PipeWriter
	ch chan readFromDone
}

func MakePipedReaderFrom(r io.ReaderFrom) PipedReader {
	var p pipedReaderFrom

	var pr *io.PipeReader
	pr, p.PipeWriter = io.Pipe()
	p.ch = make(chan readFromDone)

	go func() {
		var msg readFromDone

		if msg.n, msg.err = r.ReadFrom(pr); msg.err != nil {
			if !errors.IsEOF(msg.err) {
				pr.CloseWithError(msg.err)
			}
		}

		p.ch <- msg
	}()

	return p
}

func (p pipedReaderFrom) Close() (n int64, err error) {
	if p.PipeWriter == nil {
		return
	}

	p.PipeWriter.Close()
	out := <-p.ch
	n = out.n
	err = out.err

	return
}

type pipedDecoderFrom struct {
	*io.PipeWriter
	ch chan readFromDone
}

func MakePipedDecoder[O any](
	object O,
	decoder interfaces.DecoderFrom[O],
) PipedReader {
	var p pipedDecoderFrom

	var pr *io.PipeReader
	pr, p.PipeWriter = io.Pipe()
	p.ch = make(chan readFromDone)

	go func() {
		var msg readFromDone

		if msg.n, msg.err = decoder.DecodeFrom(object, pr); msg.err != nil {
			if !errors.IsEOF(msg.err) {
				pr.CloseWithError(msg.err)
			}
		}

		p.ch <- msg
	}()

	return p
}

func (p pipedDecoderFrom) Close() (n int64, err error) {
	if p.PipeWriter == nil {
		return
	}

	p.PipeWriter.Close()
	out := <-p.ch
	n = out.n
	err = out.err

	return
}
