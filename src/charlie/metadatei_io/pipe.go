package metadatei_io

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type pipedReaderFrom struct {
	*io.PipeWriter
	ch chan readFromDone
}

type readFromDone struct {
	n   int64
	err error
}

func makePipedReaderFrom(r io.ReaderFrom) (p pipedReaderFrom) {
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

	return
}

func (p pipedReaderFrom) Close() (out readFromDone) {
	if p.PipeWriter != nil {
		p.PipeWriter.Close()
		out = <-p.ch
	}

	return
}
