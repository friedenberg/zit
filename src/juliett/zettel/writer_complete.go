package zettel

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type WriterComplete struct {
	wBuf         *bufio.Writer
	chTransacted chan Transacted
	chDone       chan struct{}
}

func MakeWriterComplete(w io.Writer) WriterComplete {
	w1 := WriterComplete{
		chTransacted: make(chan Transacted),
		chDone:       make(chan struct{}),
		wBuf:         bufio.NewWriter(w),
	}

	go func(s *WriterComplete) {
		for z := range s.chTransacted {
			if z.Sku.Kennung.String() == "/" {
				errors.Err().Printf("empty: %#v", z)
				continue
			}

			errors.TodoP4("handle write errors")
			s.wBuf.WriteString(z.Sku.Kennung.String())
			s.wBuf.WriteString("\tZettel: !")
			s.wBuf.WriteString(z.Objekte.Typ.String())
			s.wBuf.WriteString(" ")
			s.wBuf.WriteString(z.Objekte.Bezeichnung.String())
			s.wBuf.WriteString("\n")
		}

		s.chDone <- struct{}{}
	}(&w1)

	return w1
}

func (w *WriterComplete) WriteZettelVerzeichnisse(z *Transacted) (err error) {
	select {
	case <-w.chDone:
		err = collections.ErrStopIteration

	case w.chTransacted <- *z:
	}

	return
}

func (w *WriterComplete) Close() (err error) {
	close(w.chTransacted)
	<-w.chDone

	if err = w.wBuf.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}