package zettel_named

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type WriterComplete struct {
	wBuf     *bufio.Writer
	chZettel chan Zettel
	chDone   chan struct{}
}

func MakeWriterComplete(w io.Writer) WriterComplete {
	w1 := WriterComplete{
		chZettel: make(chan Zettel),
		chDone:   make(chan struct{}),
		wBuf:     bufio.NewWriter(w),
	}

	go func(s *WriterComplete) {
		for z := range s.chZettel {
			//TODO handle errors
			s.wBuf.WriteString(z.Hinweis.String())
			s.wBuf.WriteString("\tZettel: !")
			s.wBuf.WriteString(z.Stored.Zettel.Typ.String())
			s.wBuf.WriteString(" ")
			s.wBuf.WriteString(z.Stored.Zettel.Bezeichnung.String())
			s.wBuf.WriteString("\n")
		}

		s.chDone <- struct{}{}
	}(&w1)

	return w1
}

func (w *WriterComplete) WriteZettelNamed(z Zettel) (err error) {
	select {
	case <-w.chDone:
		err = io.EOF

	case w.chZettel <- z:
	}

	return
}

func (w *WriterComplete) Close() (err error) {
	close(w.chZettel)
	<-w.chDone

	if err = w.wBuf.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
