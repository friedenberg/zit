package zettel

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type WriterComplete struct {
	wBuf    *bufio.Writer
	chNamed chan Named
	chDone  chan struct{}
}

func MakeWriterComplete(w io.Writer) WriterComplete {
	w1 := WriterComplete{
		chNamed: make(chan Named),
		chDone:  make(chan struct{}),
		wBuf:    bufio.NewWriter(w),
	}

	go func(s *WriterComplete) {
		for z := range s.chNamed {
			if z.Kennung.String() == "/" {
				errors.Err().Printf("empty: %#v", z)
				continue
			}

			//TODO-P4 handle write errors
			s.wBuf.WriteString(z.Kennung.String())
			s.wBuf.WriteString("\tZettel: !")
			s.wBuf.WriteString(z.Stored.Objekte.Typ.String())
			s.wBuf.WriteString(" ")
			s.wBuf.WriteString(z.Stored.Objekte.Bezeichnung.String())
			s.wBuf.WriteString("\n")
		}

		s.chDone <- struct{}{}
	}(&w1)

	return w1
}

func (w *WriterComplete) WriteZettelNamed(z *Named) (err error) {
	select {
	case <-w.chDone:
		err = io.EOF

	case w.chNamed <- *z:
	}

	return
}

func (w *WriterComplete) Close() (err error) {
	close(w.chNamed)
	<-w.chDone

	if err = w.wBuf.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
