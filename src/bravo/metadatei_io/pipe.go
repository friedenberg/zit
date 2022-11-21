package metadatei_io

import "io"

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
	p.ch = make(chan readFromDone, 1)

	go func() {
		var msg readFromDone
		msg.n, msg.err = r.ReadFrom(pr)
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
